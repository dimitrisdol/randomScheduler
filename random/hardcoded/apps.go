// Package hardcoded contains an implementation of random.InterferenceModel,
// where the slowdowns among all applications are known and hardcoded.
package hardcoded

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

///////////////////////////////////////////////////////////////////////////////
//
// HardcodedSlowDowns
//
///////////////////////////////////////////////////////////////////////////////

// HardcodedSlowDowns is an implementation of random.InterferenceModel, where the
// slowdowns among all applications are known and hardcoded.
type HardcodedSlowDowns struct {
	randomLabelKey string
	categoryLabelKey string
}

// New returns a new HardcodedSlowDowns with the given label key (the one that
// is used by RandomPlugin to track its applications).
func New(randomLabelKey string) *HardcodedSlowDowns {
	return &HardcodedSlowDowns{
		randomLabelKey: randomLabelKey,
	}
}

// Attack implements random.InterferenceModel; see the documentation there for
// more information.
func (m *HardcodedSlowDowns) Attack(attacker, occupant *corev1.Pod) (float64, error) {
	occPodCategory, _ := parseAppCategory(occupant.Labels[m.categoryLabelKey])
	//  occupant   ^^^   Pod's label's value must have been
	// validated back when it was first scheduled on the Node
	newPodCategory, err := parseAppCategory(attacker.Labels[m.categoryLabelKey])
	if err != nil {
		return -1, err
	}
	return newPodCategory.attack(occPodCategory), nil
}

const toInt64Multiplier = 100.

// ToInt64Multiplier implements random.InterferenceModel; see the documentation
// there for more information.
func (_ *HardcodedSlowDowns) ToInt64Multiplier() float64 {
	return toInt64Multiplier
}

///////////////////////////////////////////////////////////////////////////////
//
// appCategory
//
///////////////////////////////////////////////////////////////////////////////

// appCategory is an enumeration of known application categories.
// the 4 appCategories that are to be examined by the RandomPlugin are as
// follows:
//  catA : insensitive and peaceful applications, considered the best for the
//         plugin
//  catB : sensitive and peaceful applications, can fit with other peaceful only
//  catC : insensitive and aggressive applications, can fit with itself and catA
//  catD : sensitive and aggressive applications, the worst kind and can only
//         fit with catA
//
//  All of these are implemented by the slowdown matrix later

type appCategory int64

const (
	catA appCategory = iota
	catB
	catC
	catD
)

// String returns the string representation of the (known) appCategory.
func (ac appCategory) String() string {
	switch ac {
	case catA:
		return "catA"
	case catB:
		return "catB"
	case catC:
		return "catC"
	case catD:
		return "catD"
	default:
		return "UNKNOWN"
	}
}

// parseAppCategory parses a (known) appCategory from a string.
func parseAppCategory(category string) (appCategory, error) {
	switch category {
	case "catA":
		return catA, nil
	case "catB":
		return catB, nil
	case "catC":
		return catC, nil
	case "catD":
		return catD, nil
	default:
		return -1, fmt.Errorf("unknown application category: '%s'", category)
	}
}

// attack returns the slowdown incurred on the given occupant when the
// appCategory is scheduled along with it.
func (ac appCategory) attack(occupant appCategory) float64 {
	return slowDowns[ac][occupant]
}

///////////////////////////////////////////////////////////////////////////////
//
// slowDownMatrix
//
///////////////////////////////////////////////////////////////////////////////

// slowDownMatrix is a type alias for internal use in RandomPlugin.
// Because it is a random algorithm every slowdown is the same.
type slowDownMatrix map[appCategory]map[appCategory]float64

// slowDowns is a hardcoded global map that represents a dense 2D matrix of the
// slowdowns incurred by application colocations. Its format is as follows:
//
//     {
//         A: {
//             A: f64 slowdown of an A when attacked by an A
//             B: f64 slowdown of a B when attacked by an A
//             C: f64 slowdown of a C when attacked by an A
//         },
//         B: {
//             A: f64 slowdown of an A when attacked by a B
//             B: f64 slowdown of a B when attacked by a B
//             C: f64 slowdown of a C when attacked by a B
//         },
//         C: {
//             A: f64 slowdown of an A when attacked by a C
//             B: f64 slowdown of a B when attacked by a C
//             C: f64 slowdown of a C when attacked by a C
//         },
//         . . .
//     }
var slowDowns = slowDownMatrix{
	catA: map[appCategory]float64{
		catA: 1.00,
		catB: 1.00, // slowdown of catB when attacked by catA = 1.00
		catC: 1.00,
		catD: 1.00,
	},
	catB: map[appCategory]float64{
		catA: 1.00,
		catB: 1.00,
		catC: 1.00,
		catD: 1.00, // slowdown of catD when attacked by catB = 1.00
	},
	catC: map[appCategory]float64{
		catA: 1.00,
		catB: 1.00,
		catC: 1.00, // slowdown of catC when attacked by catC = 1.00
		catD: 1.00,
	},
	catD: map[appCategory]float64{
		catA: 1.00, // slowdown of catA when attacked by catD = 1.00
		catB: 1.00,
		catC: 1.00,
		catD: 1.00,
	},

}

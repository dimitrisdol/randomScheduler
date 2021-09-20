package greedy

import corev1 "k8s.io/api/core/v1"

// InterferenceModel abstracts away the core logic of interference awareness.
//
// It is used by RandomPlugin across various scheduling extension points, and
// exposes the generalized functionality that should be provided to RandomPlugin.
type InterferenceModel interface {
	// Attack returns a float64 that represents the damage inflicted to the
	// occupant Pod when the attacker Pod is scheduled along with it.
	Attack(attacker, occupant *corev1.Pod) (float64, error)

	// ToInt64Multiplier returns an integer that can be used to safely
	// convert the result of Attack into an int64.
	//
	// The aim is to multiply the double-precision floating point number
	// returned by Attack with this number before casting it to an int64.
	ToInt64Multiplier() float64
}

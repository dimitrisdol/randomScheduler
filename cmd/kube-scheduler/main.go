package main

import (
	"github.com/dimitrisdol/randomScheduler/random"

	"k8s.io/klog/v2"
	sched "k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	cmd := sched.NewSchedulerCommand(
		sched.WithPlugin(random.Name, random.New),
	)
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("failed to execute %q: %v", random.Name, err)
	}
}

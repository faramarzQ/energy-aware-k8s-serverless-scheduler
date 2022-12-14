package main

import (
	"fmt"

	cluster "energy-aware-k8s-serverless-scheduler/pkg/cluster"
)

func main() {
	fmt.Println("--------------------------")
	fmt.Println("Running Energy Adapter")
	fmt.Println("--------------------------")

	runFunctionOffloader()
}

func runFunctionOffloader() {
	// list all pods
	pods := cluster.ListPods()

	for _, pod := range pods {
		exists := ifABetterNodeExistForPod(pod)

		if exists {
			// delete the pod so that the Deployment creates it again (offloading)
			pod.Delete()
			fmt.Println("Offloading pod: ", pod.Name)
		}
	}
}

func ifABetterNodeExistForPod(pod cluster.Pod) bool {
	var exists bool

	nodes := cluster.ListNodes()
	podsNode := pod.GetNode()

	for _, node := range nodes {

		if podsNode.IfEnergyStatusIsLowerThan(node.EnergyStatus()) {
			exists = true
		}
	}

	return exists
}

package main

import (
	cluster "energy-aware-k8s-serverless-scheduler/pkg/cluster"
	"fmt"
)

func main() {
	fmt.Println("--------------------------")
	fmt.Println("Running Energy Adapter")
	fmt.Println("--------------------------")

	bindNodesEnergyStatus()
}

// Bind new energy status to each node
func bindNodesEnergyStatus() {
	nodes := cluster.ListNodes()

	if nodes.IsEnergyStatusInitialized() {
		nodes.UpdateEnergyStatus()
		return
	}

	nodes.InitializeNodesEnergyStatus()
}

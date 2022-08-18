package main

import (
	cluster "energy-aware-k8s-serverless-scheduler/pkg/cluster"
)

func main() {
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

package main

import (
	"context"
	cluster "energy-aware-k8s-serverless-scheduler/pkg/cluster"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type EnergyAwareScheduler struct {
	framework.FilterPlugin
	// framework.ScorePlugin
}

const Name = "EnergyAwareScheduler"

func (pl *EnergyAwareScheduler) Name() string {
	return Name
}

func (cs *EnergyAwareScheduler) Filter(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	pod := cluster.Pod{*p}
	node := (cluster.Node{*nodeInfo.Node()})
	var status framework.Code = framework.Success

	// check energy level
	if node.IfNodeIsPowerLess() {
		status = framework.Unschedulable
	}

	// check number of nodes
	if node.IfNodeCapacityIsFull() {
		klog.Info("Capacity is full on node ", node.Name)
		status = framework.Unschedulable
	}

	if ifANodeWithHigherEnergyStatusExists(node) {
		status = framework.Unschedulable
	}

	klog.Info("Filtering Pod: ", pod.Name, " on ", nodeInfo.Node().Name, " : ", status.String())

	return framework.NewStatus(status)
}

func ifANodeWithHigherEnergyStatusExists(targetNode cluster.Node) bool {
	var exists bool

	for _, node := range cluster.ListNodes() {
		if node.Name == targetNode.Name {
			continue
		}

		if targetNode.IfEnergyStatusIsLowerThan(node.EnergyStatus()) {
			exists = true
		}

		if node.IfNodeCapacityIsFull() {
			exists = false
		}

		if exists == true {
			break
		}

	}

	return exists
}

func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &EnergyAwareScheduler{}, nil
}

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(Name, New),
	)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

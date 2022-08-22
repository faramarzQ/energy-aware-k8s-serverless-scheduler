package cluster

import (
	"context"
	"fmt"

	"energy-aware-k8s-serverless-scheduler/pkg/clientSet"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/klog"
)

type Node struct {
	v1.Node
}

func (n *Node) EnergyStatus() string {
	return n.Labels["energy-status"]
}

func (n *Node) Delete() string {
	return n.Delete()
}

// if node's energy status is lower than the given node's energy
func (n *Node) IfEnergyStatusIsLowerThan(energyStatus string) bool {
	var isHigher bool

	if n.EnergyStatus() == "well-powered" {
		isHigher = false
	} else if n.EnergyStatus() == "vulnerable" {
		if energyStatus == "well-powered" {
			isHigher = true
		}
	} else if n.EnergyStatus() == "low-powered" {
		if energyStatus == "well-powered" ||
			energyStatus == "vulnerable" {
			isHigher = true
		}
	} else if n.EnergyStatus() == "powerless" {
		if energyStatus == "well-powered" ||
			energyStatus == "vulnerable" ||
			energyStatus == "low-powered" {
			isHigher = true
		}
	}

	return isHigher
}

func (n Node) SetEnergyStatus(status string) {
	c := clientSet.GetClientset()

	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, "energy-status", status)
	_, err := c.CoreV1().Nodes().Patch(context.Background(), n.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

func (n Node) ListPods() PodList {
	c := clientSet.GetClientset()

	pods, err := c.CoreV1().Pods("openfaas-fn").List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + n.Name,
	})
	if err != nil {
		panic(err)
	}

	podList := PodList{}
	for _, node := range pods.Items {
		podList = append(podList, Pod{node})
	}

	return podList
}

func (n Node) IfNodeCapacityIsFull() bool {
	if len(n.ListPods()) == 4 {
		return true
	}
	return false
}

func (n Node) IfNodeIsPowerLess() bool {
	if n.EnergyStatus() == "powerless" {
		return true
	}
	return false
}

type NodeList []Node

// gets list if nodes from cluster
func listNodesFomCluster() *v1.NodeList {
	c := clientSet.GetClientset()
	nodes, err := c.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	return nodes
}

func ListNodes() NodeList {
	nodes := listNodesFomCluster()

	nodeList := NodeList{}
	for _, node := range nodes.Items {
		nodeList = append(nodeList, Node{node})
	}

	return nodeList
}

func (nl NodeList) IsEnergyStatusInitialized() bool {
	// if there is one node with empty energy status
	for _, node := range nl {
		if node.Labels["energy-status"] == "" {
			return false
		}
	}
	return true
}

func (nl NodeList) UpdateEnergyStatus() {
	// sets new energy status for each node
	fmt.Println("Updating nodes energy status")

	updateStatuses := map[string]string{
		"well-powered": "low-powered",
		"vulnerable":   "powerless",
		"low-powered":  "well-powered",
		"powerless":    "vulnerable",
	}

	for _, node := range nl {
		currentStatus := node.EnergyStatus()
		node.SetEnergyStatus(updateStatuses[currentStatus])
		fmt.Println(node.Name, currentStatus, " => ", updateStatuses[currentStatus])
	}
}

// initialize new energy status for each node
func (nl NodeList) InitializeNodesEnergyStatus() {
	fmt.Println("Initializing nodes energy status")

	initStatuses := []string{"well-powered", "vulnerable", "low-powered", "powerless"}
	for i, status := range initStatuses {
		nl[i].SetEnergyStatus(status)
		fmt.Println(nl[i].Name, status)
	}
}

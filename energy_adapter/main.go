package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	// v1 "k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	bindNodesEnergyStatus()
}

func getClientset() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", "../minikube_auth/config")
	if err != nil {
		fmt.Println("error on building config")
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}

// gets list if nodes
func getNodes() *v1.NodeList {
	c := getClientset()

	nodes, err := c.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	return nodes
}

// updates nodes' energy status randomly
func bindNodesEnergyStatus() {
	nodes := getNodes()

	if ifNodesEnergyStatusInitialized(nodes) {
		updateNodesEnergyStatus(nodes)
		return
	}

	initializeNodesEnergyStatus(nodes)
}

func setNodeEnergyStatus(nodeName string, status string) {
	c := getClientset()

	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, "energy-status", status)
	_, err := c.CoreV1().Nodes().Patch(context.Background(), nodeName, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

// If nodes' energy status are initialized before
func ifNodesEnergyStatusInitialized(nodes *v1.NodeList) bool {
	// if there is one node with empty energy status
	for _, node := range nodes.Items {
		if node.Labels["energy-status"] == "" {
			return false
		}
	}
	return true
}

// initialize new energy status for each node
func initializeNodesEnergyStatus(nodes *v1.NodeList) {
	fmt.Println("Initializing nodes energy status")

	initStatuses := []string{"well-powered", "vulnerable", "low-powered", "powerless"}
	for i, status := range initStatuses {
		setNodeEnergyStatus(nodes.Items[i].Name, status)
		fmt.Println(nodes.Items[i].Name, status)

	}
}

// sets new energy status for each node
func updateNodesEnergyStatus(nodes *v1.NodeList) {
	fmt.Println("Updating nodes energy status")

	updateStatuses := map[string]string{
		"well-powered": "low-powered",
		"vulnerable":   "powerless",
		"low-powered":  "well-powered",
		"powerless":    "vulnerable",
	}

	for _, node := range nodes.Items {
		currentStatus := node.Labels["energy-status"]
		setNodeEnergyStatus(node.Name, updateStatuses[currentStatus])
		fmt.Println(node.Name, currentStatus, " => ", updateStatuses[currentStatus])
	}
}

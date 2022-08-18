package cluster

import (
	"context"
	"energy-aware-k8s-serverless-scheduler/pkg/clientSet"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type Pod struct {
	v1.Pod
}

type PodList []Pod

// gets list if nodes from cluster
func getPodsFomCluster() *v1.PodList {
	c := clientSet.GetClientset()
	pods, err := c.CoreV1().Pods("openfaas-fn").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	return pods
}

func ListPods() PodList {
	pods := getPodsFomCluster()
	podList := PodList{}
	for _, node := range pods.Items {
		podList = append(podList, Pod{node})
	}

	return podList
}

// gets the node of the pod
func (p Pod) GetNode() Node {
	var podsNode Node
	nodes := ListNodes()
	for _, node := range nodes {
		if node.Name == p.Spec.NodeName {
			podsNode = node
		}
	}
	return podsNode
}

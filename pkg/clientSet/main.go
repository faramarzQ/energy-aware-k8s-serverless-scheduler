package clientSet

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", "./minikube_auth/config")
	if err != nil {
		fmt.Println("error on building clientset config")
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}

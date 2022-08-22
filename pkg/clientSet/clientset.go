package clientSet

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset() *kubernetes.Clientset {
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "../minikube_auth/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", configDir)
	if err != nil {
		fmt.Println("error on building clientset config")
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset
}

package initialize

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// InitK8sClient 初始化 k8s client
// Return *kubernetes.Clientset
func InitK8sClient(kubeconfig *string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return client
}

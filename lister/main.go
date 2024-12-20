package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Usage: go run main.go --kubeconfig ~/.kube/config

func main() {
	kubeconfig := flag.String("kubeconfig", "missing kubeconfig", "location of kubeconfig")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("clientcmd failed: %v", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("clientset failed: %v", err)
		os.Exit(1)
	}

	namespace := "kube-system"

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("pod list failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Pods from %v namespace:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("deployment list failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDeployments from %v namespace:\n", namespace)
	for _, deploy := range deployments.Items {
		fmt.Println(deploy.Name)
	}
}

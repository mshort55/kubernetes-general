package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "missing kubeconfig", "location of kubeconfig")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("unable to build client from flag: %v\n", err)
		log.Printf("trying to use pod service account token...\n")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("unable to get in cluster config: %v\n", err)
		}
	}

	config.Timeout = 10 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("unable to create clientset: %v\n", err)
	}

	namespace := "kube-system"

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("pod list failed: %v\n", err)
	}

	fmt.Printf("Pods from %v namespace:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("deployment list failed: %v\n", err)
	}

	fmt.Printf("\nDeployments from %v namespace:\n", namespace)
	for _, deploy := range deployments.Items {
		fmt.Println(deploy.Name)
	}
}

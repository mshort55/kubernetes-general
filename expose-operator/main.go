package main

import (
	"flag"
	"log"
	"time"

	"k8s.io/client-go/informers"
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("unable to create clientset: %v\n", err)
	}

	ch := make(chan struct{})
	informers := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	c := newController(clientset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.run(ch)
}

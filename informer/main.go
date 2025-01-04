package main

import (
	"flag"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	informerfactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	podinformer := informerfactory.Core().V1().Pods()
	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			log.Print("add was called")
		},
		UpdateFunc: func(old, new interface{}) {
			log.Print("updated was called")
		},
		DeleteFunc: func(obj interface{}) {
			log.Print("updated was called")
		},
	})

	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)

	pod, err := podinformer.Lister().Pods("kube-system").Get("etcd-minikube")
	log.Printf("pod info: \nname: %v\nlabels: %v\nstatus: %v", pod.Name, pod.Labels, pod.Status)
}

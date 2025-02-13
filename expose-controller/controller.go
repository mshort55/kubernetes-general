package main

import (
	"context"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset         kubernetes.Interface
	deployLister      appslisters.DeploymentLister
	deployCacheSynced cache.InformerSynced
	queue             workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, deployInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:         clientset,
		deployLister:      deployInformer.Lister(),
		deployCacheSynced: deployInformer.Informer().HasSynced,
		queue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "expose"),
	}

	deployInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDel,
	})

	return c
}

func (c *controller) run(ch <-chan struct{}) {
	log.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.deployCacheSynced) {
		log.Println("waiting for cache to be synced")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		log.Println("queue is shut down")
		return false
	}

	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("error getting key from cache: %v", err)
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("error splitting namespace and name: %v", err)
		return false
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		log.Printf("error syncing deployment: %v", err)
		return false
	}
	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err:= c.deployLister.Deployments(ns).Get(name)
	if err != nil {
		log.Printf("error getting deployment from lister: %v", err)
	}

	labels := depLabels(dep)
	
	// create service
	// TODO: figure out which port and port name container is listening on
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dep.Name,
			Namespace: dep.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	_, err = c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		log.Printf("error creating service: %v", err)
	}
	
	// create ingress
	return createIngress(ctx, c.clientset, svc)
}

func depLabels(dep *appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}

func createIngress(ctx context.Context, client kubernetes.Interface, svc corev1.Service) error {
	pathType := netv1.PathTypePrefix
	ingress := netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path: fmt.Sprintf("/%s", svc.Name),
									PathType: (*netv1.PathType)(&pathType),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: svc.Name,
											Port: netv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := client.NetworkingV1().Ingresses(svc.Namespace).Create(ctx, &ingress, metav1.CreateOptions{})
	
	return err
}


func (c *controller) handleAdd(obj interface{}) {
	log.Println("add was called")
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	log.Println("delete was called")
	c.queue.Add(obj)
}

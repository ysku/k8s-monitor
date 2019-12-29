package main

import (
	"fmt"
	"log"
	"time"

	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type PodLoggingController struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

func NewPodLoggingController(informerFactory informers.SharedInformerFactory) *PodLoggingController {
	podInformer := informerFactory.Core().V1().Pods()

	c := &PodLoggingController{
		informerFactory: informerFactory,
		podInformer: podInformer,
	}
	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}

func (c *PodLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *PodLoggingController) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("[podAdd] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

func (c *PodLoggingController) podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	fmt.Printf("[podUpdate] old, namespace:%s, name:%s, labels:%v", oldPod.Namespace, oldPod.Name, oldPod.GetLabels())
	fmt.Printf("[podUpdate] new, namespace:%s, name:%s, labels:%v", newPod.Namespace, newPod.Name, newPod.GetLabels())
}

func (c *PodLoggingController) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("[podDelete] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

func main() {
	log.Print("Shared Informer app started")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 10 * time.Second)
	controller := NewPodLoggingController(factory)
	stop := make(chan struct{})
	defer close(stop)
	err = controller.Run(stop)
	if err != nil {
		log.Printf("error : %v", err)
	}
	select {}
}

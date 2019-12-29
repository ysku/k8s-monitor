package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type PodLoggingController struct {
	informer cache.SharedInformer
	queue    workqueue.RateLimitingInterface
}

func NewPodLoggingController(queue workqueue.RateLimitingInterface, informer cache.SharedInformer) *PodLoggingController {
	c := &PodLoggingController{
		informer: informer,
		queue:    queue,
	}
	informer.AddEventHandler(
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

func (c *PodLoggingController) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

func (c *PodLoggingController) syncToStdout(key string) error {
	_, exists, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		log.Printf("error fetching object from index for the specified key : %s", key)
		return err
	}

	if !exists {
		log.Printf("pod has gone, key: %s", key)
		// do your heavy stuff for when the pod is gone here
	} else {
		log.Printf("update received for pod, key: %s", key)
		// do your heavy stuff for when the pod is created/updated here
	}
	return nil
}

func (c *PodLoggingController) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		log.Printf("error during sync, key: %s, error: %v", key.(string), err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	log.Printf("drop pod out of queue after many retries, key: %s, error: %v", key.(string), err)
}

func (c *PodLoggingController) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	fmt.Print("[Run] start controller!!")

	// Starts all the shared informers that have been created by the factory so far.
	go c.informer.Run(stopCh)

	fmt.Print("[Run] go c.informer.Run(stopCh)")
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		log.Print("Failed to sync")
		return
	}
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	log.Print("stopping controller")
}

func (c *PodLoggingController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *PodLoggingController) podAdd(obj interface{}) {
	fmt.Print("[podAdd] called!!")
	pod := obj.(*v1.Pod)
	fmt.Printf("[podAdd] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

func (c *PodLoggingController) podUpdate(old, new interface{}) {
	fmt.Print("[podUpdate] called!!")
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	fmt.Printf("[podUpdate] old, namespace:%s, name:%s, labels:%v", oldPod.Namespace, oldPod.Name, oldPod.GetLabels())
	fmt.Printf("[podUpdate] new, namespace:%s, name:%s, labels:%v", newPod.Namespace, newPod.Name, newPod.GetLabels())
}

func (c *PodLoggingController) podDelete(obj interface{}) {
	fmt.Print("[podDelete] called!!")
	pod := obj.(*v1.Pod)
	fmt.Printf("[podDelete] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

func main() {
	log.Print("[main] Shared Informer app started")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Print("[main] configured clientset")

	factory := informers.NewSharedInformerFactory(clientset, 10*time.Second)
	informer := factory.Core().V1().Pods().Informer()
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	log.Print("[main] configured factory")
	controller := NewPodLoggingController(queue, informer)
	log.Print("[main] configured controller")
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	log.Print("[main] start running controller")
	select {}
	log.Print("[main] here should not be reached")
}

package controller

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type LoggingController struct {
	informer cache.SharedInformer
	queue    workqueue.RateLimitingInterface
}

func NewLoggingController(
	queue workqueue.RateLimitingInterface,
	informer cache.SharedInformer) *LoggingController {
	return &LoggingController{
		informer: informer,
		queue:    queue,
	}
}

func (c *LoggingController) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	log.Print("[Run] start logging!!")

	// Starts all the shared informers that have been created by the factory so far.
	go c.informer.Run(stopCh)

	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		log.Print("Failed to sync")
		return
	}
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	log.Print("stopping logging")
}

func (c *LoggingController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *LoggingController) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

func (c *LoggingController) syncToStdout(key string) error {
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

func (c *LoggingController) handleErr(err error, key interface{}) {
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

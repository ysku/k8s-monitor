package controller

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

func NewPodLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podAdd,
			UpdateFunc: podUpdate,
			DeleteFunc: podDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController(queue, informer)
}

func podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("[podAdd] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

func podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	log.Printf("[podUpdate] old, namespace:%s, name:%s, labels:%v", oldPod.Namespace, oldPod.Name, oldPod.GetLabels())
	log.Printf("[podUpdate] new, namespace:%s, name:%s, labels:%v", newPod.Namespace, newPod.Name, newPod.GetLabels())
}

func podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("[podDelete] namespace:%s, name:%s, labels:%v", pod.Namespace, pod.Name, pod.GetLabels())
}

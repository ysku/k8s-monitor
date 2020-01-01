package controller

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

func NewServiceLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Core().V1().Services().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    serviceAdd,
			UpdateFunc: serviceUpdate,
			DeleteFunc: serviceDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController(queue, informer)
}

func serviceAdd(obj interface{}) {
	service := obj.(*v1.Service)
	log.Printf("[serviceAdd] namespace:%s, name:%s, labels:%v", service.Namespace, service.Name, service.GetLabels())
}

func serviceUpdate(old, new interface{}) {
	oldService := old.(*v1.Service)
	newService := new.(*v1.Service)
	log.Printf("[serviceUpdate] old, namespace:%s, name:%s, labels:%v", oldService.Namespace, oldService.Name, oldService.GetLabels())
	log.Printf("[serviceUpdate] new, namespace:%s, name:%s, labels:%v", newService.Namespace, newService.Name, newService.GetLabels())
}

func serviceDelete(obj interface{}) {
	service := obj.(*v1.Service)
	log.Printf("[serviceDelete] namespace:%s, name:%s, labels:%v", service.Namespace, service.Name, service.GetLabels())
}

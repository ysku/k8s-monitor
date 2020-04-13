package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
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
	return NewLoggingController("service", queue, informer)
}

func logService(action string, service *v1.Service) {
	log.WithFields(log.Fields{
		"namespace": service.Namespace,
		"name": service.Name,
		"labels": service.GetLabels(),
		"annotations": service.GetAnnotations(),
	}).Info(action)
}

func serviceAdd(obj interface{}) {
	service := obj.(*v1.Service)
	logService("[serviceAdd]", service)
}

func serviceUpdate(old, new interface{}) {
	oldService := old.(*v1.Service)
	newService := new.(*v1.Service)
	logService("[serviceUpdate] old", oldService)
	logService("[serviceUpdate] new", newService)
}

func serviceDelete(obj interface{}) {
	service := obj.(*v1.Service)
	logService("[serviceDelete]", service)
}

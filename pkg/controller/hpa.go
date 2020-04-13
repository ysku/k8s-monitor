package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/autoscaling/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewHPALoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Autoscaling().V1().HorizontalPodAutoscalers().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    hpaAdd,
			UpdateFunc: hpaUpdate,
			DeleteFunc: hpaDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController("hpa", queue, informer)
}

func hpaAdd(obj interface{}) {
	hpa := obj.(*v1.HorizontalPodAutoscaler)
	log.Printf("[hpaAdd] namespace:%s, name:%s, labels:%v", hpa.Namespace, hpa.Name, hpa.GetLabels())
}

func hpaUpdate(old, new interface{}) {
	oldHpa := old.(*v1.HorizontalPodAutoscaler)
	newHpa := new.(*v1.HorizontalPodAutoscaler)
	log.Printf("[hpaUpdate] old, namespace:%s, name:%s, labels:%v", oldHpa.Namespace, oldHpa.Name, oldHpa.GetLabels())
	log.Printf("[hpaUpdate] new, namespace:%s, name:%s, labels:%v", newHpa.Namespace, newHpa.Name, newHpa.GetLabels())
}

func hpaDelete(obj interface{}) {
	hpa := obj.(*v1.HorizontalPodAutoscaler)
	log.Printf("[hpaDelete] namespace:%s, name:%s, labels:%v", hpa.Namespace, hpa.Name, hpa.GetLabels())
}

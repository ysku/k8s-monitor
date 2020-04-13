package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/ysku/my-k8s-custom-controller/pkg/logging"
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
	return NewLoggingController("pod", queue, informer)
}

func podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.WithFields(logging.MapToFields(pod.GetLabels())).Infof("[podAdd] namespace:%s, name:%s", pod.Namespace, pod.Name)
}

func podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	log.WithFields(logging.MapToFields(oldPod.GetLabels())).Infof("[podUpdate] old, namespace:%s, name:%s", oldPod.Namespace, oldPod.Name)
	log.WithFields(logging.MapToFields(newPod.GetLabels())).Infof("[podUpdate] new, namespace:%s, name:%s", newPod.Namespace, newPod.Name)
}

func podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.WithFields(logging.MapToFields(pod.GetLabels())).Infof("[podDelete] namespace:%s, name:%s", pod.Namespace, pod.Name)
}

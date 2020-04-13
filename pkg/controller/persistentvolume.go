package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewPersistentVolumeLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Core().V1().PersistentVolumes().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    persistentPersistentVolumeAdd,
			UpdateFunc: persistentPersistentVolumeUpdate,
			DeleteFunc: persistentPersistentVolumeDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController("pv", queue, informer)
}

func persistentPersistentVolumeAdd(obj interface{}) {
	persistentVolume := obj.(*v1.PersistentVolume)
	log.Printf("[persistentVolumeAdd] namespace:%s, name:%s, labels:%v", persistentVolume.Namespace, persistentVolume.Name, persistentVolume.GetLabels())
}

func persistentPersistentVolumeUpdate(old, new interface{}) {
	oldPersistentVolume := old.(*v1.PersistentVolume)
	newPersistentVolume := new.(*v1.PersistentVolume)
	log.Printf("[persistentPersistentVolumeUpdate] old, namespace:%s, name:%s, labels:%v", oldPersistentVolume.Namespace, oldPersistentVolume.Name, oldPersistentVolume.GetLabels())
	log.Printf("[persistentPersistentVolumeUpdate] new, namespace:%s, name:%s, labels:%v", newPersistentVolume.Namespace, newPersistentVolume.Name, newPersistentVolume.GetLabels())
}

func persistentPersistentVolumeDelete(obj interface{}) {
	persistentPersistentVolume := obj.(*v1.PersistentVolume)
	log.Printf("[persistentPersistentVolumeDelete] namespace:%s, name:%s, labels:%v", persistentPersistentVolume.Namespace, persistentPersistentVolume.Name, persistentPersistentVolume.GetLabels())
}

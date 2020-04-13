package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/policy/v1beta1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewPDBLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Policy().V1beta1().PodDisruptionBudgets().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    pdbAdd,
			UpdateFunc: pdbUpdate,
			DeleteFunc: pdbDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController("pdb", queue, informer)
}

func pdbAdd(obj interface{}) {
	pdb := obj.(*v1.PodDisruptionBudget)
	log.Printf("[pdbAdd] namespace:%s, name:%s, labels:%v", pdb.Namespace, pdb.Name, pdb.GetLabels())
}

func pdbUpdate(old, new interface{}) {
	oldPdb := old.(*v1.PodDisruptionBudget)
	newPdb := new.(*v1.PodDisruptionBudget)
	log.Printf("[pdbUpdate] old, namespace:%s, name:%s, labels:%v", oldPdb.Namespace, oldPdb.Name, oldPdb.GetLabels())
	log.Printf("[pdbUpdate] new, namespace:%s, name:%s, labels:%v", newPdb.Namespace, newPdb.Name, newPdb.GetLabels())
}

func pdbDelete(obj interface{}) {
	pdb := obj.(*v1.PodDisruptionBudget)
	log.Printf("[pdbDelete] namespace:%s, name:%s, labels:%v", pdb.Namespace, pdb.Name, pdb.GetLabels())
}

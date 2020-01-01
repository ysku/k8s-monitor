package controller

import (
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

func NewDeploymentLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    deploymentAdd,
			UpdateFunc: deploymentUpdate,
			DeleteFunc: deploymentDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController(queue, informer)
}

func deploymentAdd(obj interface{}) {
	deployment := obj.(*v1.Deployment)
	log.Printf("[deploymentAdd] namespace:%s, name:%s, labels:%v", deployment.Namespace, deployment.Name, deployment.GetLabels())
}

func deploymentUpdate(old, new interface{}) {
	oldDeployment := old.(*v1.Deployment)
	newDeployment := new.(*v1.Deployment)
	log.Printf("[deploymentUpdate] old, namespace:%s, name:%s, labels:%v", oldDeployment.Namespace, oldDeployment.Name, oldDeployment.GetLabels())
	log.Printf("[deploymentUpdate] new, namespace:%s, name:%s, labels:%v", newDeployment.Namespace, newDeployment.Name, newDeployment.GetLabels())
}

func deploymentDelete(obj interface{}) {
	deployment := obj.(*v1.Deployment)
	log.Printf("[deploymentDelete] namespace:%s, name:%s, labels:%v", deployment.Namespace, deployment.Name, deployment.GetLabels())
}

package controller

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewJobLoggingController(factory informers.SharedInformerFactory) *LoggingController {
	informer := factory.Batch().V1().Jobs().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    jobAdd,
			UpdateFunc: jobUpdate,
			DeleteFunc: jobDelete,
		},
	)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return NewLoggingController("job", queue, informer)
}

func jobAdd(obj interface{}) {
	job := obj.(*v1.Job)
	log.Printf("[jobAdd] namespace:%s, name:%s, labels:%v", job.Namespace, job.Name, job.GetLabels())
}

func jobUpdate(old, new interface{}) {
	oldJob := old.(*v1.Job)
	newJob := new.(*v1.Job)
	log.Printf("[jobUpdate] old, namespace:%s, name:%s, labels:%v", oldJob.Namespace, oldJob.Name, oldJob.GetLabels())
	log.Printf("[jobUpdate] new, namespace:%s, name:%s, labels:%v", newJob.Namespace, newJob.Name, newJob.GetLabels())
}

func jobDelete(obj interface{}) {
	job := obj.(*v1.Job)
	log.Printf("[jobDelete] namespace:%s, name:%s, labels:%v", job.Namespace, job.Name, job.GetLabels())
}

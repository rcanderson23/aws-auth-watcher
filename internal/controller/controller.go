package controller

import (
	"k8s.io/klog"
	"time"

	"github.com/rcanderson23/aws-auth-watcher/internal/notification"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Watcher struct {
	AuthConfigMap *AuthConfigMap
	Controller    cache.Controller
}

func NewWatcher(clientset *kubernetes.Clientset, acm *AuthConfigMap) *Watcher {

	authListWatch := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", "kube-system", fields.OneTermEqualSelector("metadata.name", "aws-auth"))

	_, controller := cache.NewInformer(authListWatch, &v1.ConfigMap{}, time.Second*0, cache.ResourceEventHandlerFuncs{
		AddFunc:    acm.Add,
		DeleteFunc: acm.Delete,
		UpdateFunc: acm.Update,
	})

	return &Watcher{
		AuthConfigMap: acm,
		Controller:    controller,
	}
}

type AuthConfigMap struct {
	AwsAuth *v1.ConfigMap
	AwsSns  *notification.AwsSns
}

func (a *AuthConfigMap) Add(obj interface{}) {
	klog.Info("aws-auth added to watcher")
	// Need to account for the aws-auth ConfigMap changing before after controller creation and before watcher
	if a.AwsAuth.ResourceVersion != obj.(*v1.ConfigMap).ResourceVersion {
		klog.Info("Auth has changed! Firing notification!")
		a.AwsSns.PublishChange(a.AwsAuth, obj)
	}
}

func (a *AuthConfigMap) Delete(obj interface{}) {
	klog.Info("aws-auth deleted! Firing notification!")
}

func (a *AuthConfigMap) Update(oldObj, newObj interface{}) {
	klog.Info("Auth has changed! Firing notification!")
	a.AwsSns.PublishChange(oldObj, newObj)
}

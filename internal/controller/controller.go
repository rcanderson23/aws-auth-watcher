package controller

import (
	"k8s.io/klog"
	"reflect"
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
	eq := reflect.DeepEqual(a.AwsAuth.Data, obj.(*v1.ConfigMap).Data)
	if !eq {
		klog.Info("Auth has changed! Firing notification!")
		a.AwsSns.PublishChange(a.AwsAuth, obj)

		// Necessary copy when watcher restarts
		obj.(*v1.ConfigMap).DeepCopyInto(a.AwsAuth)
	}
}

func (a *AuthConfigMap) Delete(obj interface{}) {
	klog.Info("aws-auth deleted! Firing notification!")
	a.AwsAuth = &v1.ConfigMap{}
	a.AwsSns.PublishDelete(obj)
}

func (a *AuthConfigMap) Update(oldObj, newObj interface{}) {
	eq := reflect.DeepEqual(oldObj.(*v1.ConfigMap).Data, newObj.(*v1.ConfigMap).Data)
	if !eq {
		klog.Info("Auth has changed! Firing notification!")
		a.AwsSns.PublishChange(a.AwsAuth, newObj)

		// Necessary copy when watcher restarts
		newObj.(*v1.ConfigMap).DeepCopyInto(a.AwsAuth)
	}
}

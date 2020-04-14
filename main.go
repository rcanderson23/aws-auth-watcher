package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	//	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Watcher struct {
	awsAuth       v1.ConfigMap
	authListWatch *cache.ListWatch
	controller    cache.Controller
}

func NewWatcher(clientset *kubernetes.Clientset) *Watcher {

	authListWatch := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", "kube-system", fields.OneTermEqualSelector("metadata.name", "aws-auth"))

	_, controller := cache.NewInformer(authListWatch, &v1.ConfigMap{}, time.Second*0, cache.ResourceEventHandlerFuncs{
		AddFunc:    Add,
		DeleteFunc: Delete,
		UpdateFunc: Update,
	})

	return &Watcher{
		awsAuth:       v1.ConfigMap{},
		authListWatch: authListWatch,
		controller:    controller,
	}
}

func Add(obj interface{}) {
	klog.Info("cm added")
}

func Delete(obj interface{}) {
	klog.Info("cm deleted")
}

func Update(oldObj, newObj interface{}) {
	klog.Info("cm updated")
}

func main() {
	config := createConfig()
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	watcher := NewWatcher(clientset)

	stop := make(chan struct{})
	go watcher.controller.Run(stop)
	select {}
}

func createConfig() *rest.Config {

	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		return config
	}
	klog.Info("Using in-cluster config")
	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

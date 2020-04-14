package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Watcher struct {
	AuthConfigMap AuthConfigMap
	AuthListWatch *cache.ListWatch
	Controller    cache.Controller
}

func NewWatcher(clientset *kubernetes.Clientset) *Watcher {

	authListWatch := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", "kube-system", fields.OneTermEqualSelector("metadata.name", "aws-auth"))

	var cm *v1.ConfigMap
	var err error
	cm, err = clientset.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
	if err != nil {
		cm = &v1.ConfigMap{}
	}

	authCM := AuthConfigMap{
		AwsAuth: cm,
	}
	_, controller := cache.NewInformer(authListWatch, &v1.ConfigMap{}, time.Second*0, cache.ResourceEventHandlerFuncs{
		AddFunc:    authCM.Add,
		DeleteFunc: authCM.Delete,
		UpdateFunc: authCM.Update,
	})

	return &Watcher{
		AuthConfigMap: authCM,
		AuthListWatch: authListWatch,
		Controller:    controller,
	}
}

type AuthConfigMap struct {
	AwsAuth *v1.ConfigMap
}

func (a *AuthConfigMap) Add(obj interface{}) {
	klog.Info("aws-auth added to watcher")
	// Need to account for the aws-auth ConfigMap changing before after controller creation and before watcher
	if a.AwsAuth.ResourceVersion != obj.(*v1.ConfigMap).ResourceVersion {
		klog.Info("Auth has changed! Firing notification!")
	}
}

func (a *AuthConfigMap) Delete(obj interface{}) {
	klog.Info("aws-auth deleted! Firing notification!")
}

func (a *AuthConfigMap) Update(oldObj, newObj interface{}) {
	klog.Info("Auth has changed! Firing notification!")
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
	go watcher.Controller.Run(stop)
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

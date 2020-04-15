package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/rcanderson23/aws-auth-watcher/internal"
	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config := createConfig()
	// creates the clientset
	var err error
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// seed current aws-auth ConfigMap
	var cm *v1.ConfigMap
	cm, err = clientset.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
	if err != nil {
		cm = &v1.ConfigMap{}
	}
	acm := internal.AuthConfigMap{
		AwsAuth: cm,
	}

	// create watcher object
	watcher := internal.NewWatcher(clientset, &acm)

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

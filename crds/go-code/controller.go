package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	// infcorev1 "k8s.io/client-go/informers/core/v1"
	"github.com/davecgh/go-spew/spew"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"time"
)

func startListening() {
	fmt.Printf("Entered startListening\n")

	cfg := config.GetConfigOrDie()
	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

	stopCh := signals.SetupSignalHandler()

	namespaceInformer := kubeInformerFactory.Core().V1().Namespaces()
	namespaceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    namespaceAdded,
		UpdateFunc: namespaceUpdated,
		DeleteFunc: namespaceDeleted,
	})

	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    deploymentAdded,
		UpdateFunc: deploymentUpdated,
		DeleteFunc: deploymentDeleted,
	})

	//TODO add an informer for a CRD type, either dynamically or with a generation step

	if ok := cache.WaitForCacheSync(stopCh, func() bool { return true }, func() bool { return true }); !ok {
		panic("failed to wait for caches to sync")
	}

	kubeInformerFactory.Start(stopCh)

	go wait.Until(runWorker, 10*time.Second, stopCh)

	<-stopCh

	fmt.Printf("completed do\n")
}

func runWorker() {
	fmt.Printf("Entered runWorker\n")
}

func namespaceDeleted(obj interface{}) {
	fmt.Printf("Deleted namespace was (%v)\n", spew.Sdump(obj))
}
func namespaceAdded(obj interface{}) {
	// fmt.Printf("Added namespace was (%v)\n", spew.Sdump(obj))
	if namespace, ok := obj.(*corev1.Namespace); ok {
		fmt.Printf("Added namespace named (%s)\n", namespace.Name)
	}
}

func namespaceUpdated(old, new interface{}) {
	// fmt.Printf("Updated namespace:\nold was (%v) \n new was (%v)\n", spew.Sdump(old), spew.Sdump(new))
	if oldNamespace, ok := old.(*corev1.Namespace); ok {
		if newNamespace, ok := new.(*corev1.Namespace); ok {
			fmt.Printf("Updated namespace from (%s) to (%s)\n", oldNamespace.Name, newNamespace.Name)
		}
	}
}

func deploymentDeleted(obj interface{}) {
	fmt.Printf("Deleted deployment was (%v)\n", spew.Sdump(obj))
}
func deploymentAdded(obj interface{}) {
	// fmt.Printf("Added deployment was (%v)\n", spew.Sdump(obj))
	if deployment, ok := obj.(*appsv1.Deployment); ok {
		fmt.Printf("Added deployment named (%s)\n", deployment.Name)
	}
}

func deploymentUpdated(old, new interface{}) {
	// fmt.Printf("Updated deployment:\nold was (%v) \n new was (%v)\n", spew.Sdump(old), spew.Sdump(new))
	if oldDeployment, ok := old.(*appsv1.Deployment); ok {
		if newDeployment, ok := new.(*appsv1.Deployment); ok {
			fmt.Printf("Updated deployment from (%s) to (%s)\n", oldDeployment.Name, newDeployment.Name)
		}
	}
}

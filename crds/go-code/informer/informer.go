package informer

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	// utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	// infcorev1 "k8s.io/client-go/informers/core/v1"
	"github.com/davecgh/go-spew/spew"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	// "k8s.io/client-go/kubernetes"
	kubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
	// "k8s.io/client-go/tools/cache"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	// scheme "k8s.io/sample-controller/pkg/generated/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"time"
)

const (
	namespaceName = "default"
)

var plural string

func Start(crdPlural string) {
	fmt.Printf("Entered Start\n")
	plural = crdPlural

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

	//try adding a dynamic informer
	kubeInformerFactory.InformerFor(&Xyz{}, func(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return newInformer(client, namespaceName, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, tweakListOptionsImpl)
	})

	if ok := cache.WaitForCacheSync(stopCh, func() bool { return true }, func() bool { return true }); !ok {
		panic("failed to wait for caches to sync")
	}

	kubeInformerFactory.Start(stopCh)

	go wait.Until(runWorker, 10*time.Second, stopCh)

	<-stopCh

	fmt.Printf("completed do\n")
}

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

func tweakListOptionsImpl(x *metav1.ListOptions) {}

func newInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return listImpl(client, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return watchImpl(client, options)
			},
		},
		&Xyz{},
		resyncPeriod,
		indexers,
	)
}
func setConfigDefaults(config *rest.Config) error {
	config.APIPath = "/apis"
	config.NegotiatedSerializer = codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

func getRESTClient() rest.Interface {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	err = setConfigDefaults(cfg)
	if err != nil {
		panic(err)
	}
	client, err := rest.UnversionedRESTClientFor(cfg)
	if err != nil {
		panic(err)
	}
	return client
}
func listImpl(client kubernetes.Interface, opts metav1.ListOptions) (result runtime.Object, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	fmt.Printf("options is %v\n", spew.Sdump(opts))
	res := getRESTClient().
		Get().
		Namespace(namespaceName).
		Resource(plural).
		// VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do()
	fmt.Printf("res is %v\n", spew.Sdump(res))
	err = res.Into(result)
	fmt.Printf("err is %v\n", err)
	return
}
func watchImpl(client kubernetes.Interface, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return getRESTClient().
		Get().
		Namespace(namespaceName).
		Resource(plural).
		Suffix("lyra.example.com").
		// VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

func runWorker() {
	fmt.Printf("Entered runWorker\n")
}

func namespaceDeleted(obj interface{}) {
	fmt.Printf("Deleted namespace was (%v)\n", spew.Sdump(obj))
}

func anyAdded(obj interface{}) {
	fmt.Printf("Added any was (%v)\n", spew.Sdump(obj))
}
func namespaceAdded(obj interface{}) {
	if namespace, ok := obj.(*corev1.Namespace); ok {
		fmt.Printf("Added namespace named (%s)\n", namespace.Name)
	}
}

func namespaceUpdated(old, new interface{}) {
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
	if deployment, ok := obj.(*appsv1.Deployment); ok {
		fmt.Printf("Added deployment named (%s)\n", deployment.Name)
	}
}

func deploymentUpdated(old, new interface{}) {
	if oldDeployment, ok := old.(*appsv1.Deployment); ok {
		if newDeployment, ok := new.(*appsv1.Deployment); ok {
			fmt.Printf("Updated deployment from (%s) to (%s)\n", oldDeployment.Name, newDeployment.Name)
		}
	}
}

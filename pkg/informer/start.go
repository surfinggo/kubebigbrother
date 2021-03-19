package informer

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"syscall"
)

func (i *Informer) Start() error {
	informer := informers.NewSharedInformerFactory(i.Clientset, 0)

	configMapInformer := informer.Core().V1().ConfigMaps().Informer()

	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap, ok := obj.(*v1.ConfigMap)
			if ok {
				i.ConfigMapAddFunc(configMap)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldConfigMap, ok1 := oldObj.(*v1.ConfigMap)
			newConfigMap, ok2 := newObj.(*v1.ConfigMap)
			if ok1 && ok2 {
				i.ConfigMapUpdateFunc(oldConfigMap, newConfigMap)

			}
		},
		DeleteFunc: func(obj interface{}) {
			configMap, ok := obj.(*v1.ConfigMap)
			if ok {
				i.ConfigMapDeleteFunc(configMap)
			}
		},
	})

	stopNodeNotReadyCh := make(chan struct{})
	defer close(stopNodeNotReadyCh)

	go configMapInformer.Run(stopNodeNotReadyCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	return nil
}

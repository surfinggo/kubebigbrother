package watcher

import (
	"github.com/spongeprojects/kubebigbrother/pkg/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"syscall"
)

func (app *Watcher) Start() error {
	informer := informers.NewSharedInformerFactory(app.Clientset, 0)

	configMapInformer := informer.Core().V1().ConfigMaps().Informer()

	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			configMap, ok := obj.(*v1.ConfigMap)
			if ok {
				log.Infof("created %s/%s", configMap.Namespace, configMap.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldConfigMap, ok1 := oldObj.(*v1.ConfigMap)
			newConfigMap, ok2 := newObj.(*v1.ConfigMap)
			if ok1 && ok2 {
				log.Infof("updated %s/%s", oldConfigMap.Namespace, newConfigMap.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			configMap, ok := obj.(*v1.ConfigMap)
			if ok {
				log.Infof("deleted %s/%s", configMap.Namespace, configMap.Name)
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

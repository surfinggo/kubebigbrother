package informer

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	KubeConfig string
	Resource   string

	ConfigMapAddFunc    func(configMap *v1.ConfigMap)
	ConfigMapUpdateFunc func(oldConfigMap *v1.ConfigMap, newConfigMap *v1.ConfigMap)
	ConfigMapDeleteFunc func(configMap *v1.ConfigMap)
}

type Informer struct {
	Clientset kubernetes.Interface

	ConfigMapAddFunc    func(configMap *v1.ConfigMap)
	ConfigMapUpdateFunc func(oldConfigMap *v1.ConfigMap, newConfigMap *v1.ConfigMap)
	ConfigMapDeleteFunc func(configMap *v1.ConfigMap)
}

func Setup(options Options) (*Informer, error) {
	informer := &Informer{}

	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "create kube client error")
	}

	informer.Clientset = clientset

	informer.ConfigMapAddFunc = options.ConfigMapAddFunc
	informer.ConfigMapUpdateFunc = options.ConfigMapUpdateFunc
	informer.ConfigMapDeleteFunc = options.ConfigMapDeleteFunc

	return informer, nil
}

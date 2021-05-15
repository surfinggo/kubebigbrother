package informers

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// ChannelName is name of channel
type ChannelName string

// ChannelType is type of channel
type ChannelType string

// ChannelGroup is a special type of channel, which includes a slice of channels
type ChannelGroup []ChannelName

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
	Token string `json:"token" yaml:"token"`
}

// ChannelCallback is the callback channel
type ChannelCallback struct {
	URL string `json:"url" yaml:"url"`
}

// Channel defines a channel to receive notifications
type Channel struct {
	// Name is the name of the channel
	Name ChannelName `json:"name" yaml:"name"`

	// Type is the type of the channel
	Type ChannelType `json:"type" yaml:"type"`

	Group    *ChannelGroup    `json:"group" yaml:"group"`
	Telegram *ChannelTelegram `json:"telegram" yaml:"telegram"`
	Callback *ChannelCallback `json:"callback" yaml:"callback"`
}

type Resource struct {
	// Resource is the resource to watch, e.g. "deployments.v1.apps"
	Resource string `json:"resource" yaml:"resource"`

	// NoticeWhenAdded determine whether to notice when a resource is added
	NoticeWhenAdded bool `json:"noticeWhenAdded" yaml:"noticeWhenAdded"`

	// NoticeWhenDeleted determine whether to notice when a resource is deleted
	NoticeWhenDeleted bool `json:"noticeWhenDeleted" yaml:"noticeWhenDeleted"`

	// NoticeWhenUpdated determine whether to notice when a resource is updated,
	// When UpdateOn is not nil, only notice when fields in UpdateOn is changed
	NoticeWhenUpdated bool `json:"noticeWhenUpdated" yaml:"noticeWhenUpdated"`

	// UpdateOn defines fields to watch, used with NoticeWhenUpdated
	UpdateOn []string `json:"updateOn" yaml:"updateOn"`

	// Channels defines channels to send notification
	Channels []ChannelName `json:"channels" yaml:"channels"`
}

type Namespace struct {
	// Namespace is the namespace to watch, default to "*", which means all namespaces
	Namespace string `json:"namespace" yaml:"namespace"`

	// Resources is the resources you want to watch
	Resources []Resource `json:"resources" yaml:"resources"`
}

type Config struct {
	// Namespaces defines namespaces and resources to watch
	Namespaces []Namespace `json:"namespaces" yaml:"namespaces"`

	// Channels defines channels that receive notifications
	Channels map[ChannelName]Channel `json:"channels" yaml:"channels"`

	// MinResyncPeriod is the resync period in reflectors; will be random between
	// MinResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod string `json:"minResyncPeriod" yaml:"minResyncPeriod"`
}

// LoadConfigFromFile loads config from file
func LoadConfigFromFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open error")
	}
	var config Config
	switch t := strings.ToLower(path.Ext(file)); t {
	case ".yaml":
		err = yaml.NewDecoder(f).Decode(&config)
		if err != nil {
			return nil, errors.Wrap(err, "yaml decode error")
		}
	case ".json":
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			return nil, errors.Wrap(err, "json decode error")
		}
	default:
		return nil, errors.Errorf("unsupported file type: %s", t)
	}
	return &config, nil
}

func (c Config) ResyncPeriodFunc() (func() time.Duration, error) {
	if c.MinResyncPeriod == "" {
		c.MinResyncPeriod = "12h"
	}
	duration, err := time.ParseDuration(c.MinResyncPeriod)
	if err != nil {
		return nil, errors.Wrap(err, "time.ParseDuration error")
	}
	f := float64(duration.Nanoseconds())
	// generation time.Duration between MinResyncPeriod and 2*MinResyncPeriod
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(f * factor)
	}, nil
}

type Options struct {
	KubeConfig string

	Config *Config
}

type Interface interface {
	Start(stopCh <-chan struct{})
}

type InformerSet struct {
	Factories []dynamicinformer.DynamicSharedInformerFactory
}

func (set *InformerSet) Start(stopCh <-chan struct{}) {
	for i, factory := range set.Factories {
		klog.Infof("starting factory %d/%d", i+1, len(set.Factories))
		go factory.Start(stopCh)
	}
}

func Setup(options Options) (*InformerSet, error) {
	config := options.Config

	informerSet := &InformerSet{}

	restConfig, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "get kube config error")
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "dynamic.NewForConfig error")
	}

	resyncPeriodFunc, err := config.ResyncPeriodFunc()
	if err != nil {
		return nil, errors.Wrap(err, "config.ResyncPeriodFunc error")
	}

	for i, namespace := range config.Namespaces {
		klog.Infof("setup namespace %d/%d", i+1, len(config.Namespaces))
		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
			dynamicClient, resyncPeriodFunc(), namespace.Namespace, nil)

		for j, resource := range namespace.Resources {
			klog.Infof("setup resource %d/%d, namespace: %s, resource: %s",
				j+1, len(namespace.Resources), namespace.Namespace, resource.Resource)
			gvr, _ := schema.ParseResourceArg(resource.Resource)
			if gvr == nil {
				return nil, errors.Wrapf(err, "schema.ParseResourceArg error, .Namespaces[%d].Resource[%d]: %s",
					i, j, resource.Resource)
			}
			informer := factory.ForResource(*gvr).Informer()
			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					fmt.Printf("created: %s\n", s.GetName())
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					oldS, ok1 := oldObj.(*unstructured.Unstructured)
					newS, ok2 := newObj.(*unstructured.Unstructured)
					if !ok1 || !ok2 {
						return
					}
					oldColor, ok1, err1 := unstructured.NestedString(oldS.Object, "spec", "color")
					newColor, ok2, err2 := unstructured.NestedString(newS.Object, "spec", "color")
					if !ok1 || !ok2 || err1 != nil || err2 != nil {
						fmt.Printf("updated: %s\n", newS.GetName())
					}
					fmt.Printf("updated: %s, old color: %s, new color: %s\n", newS.GetName(), oldColor, newColor)
				},
				DeleteFunc: func(obj interface{}) {
					s, ok := obj.(*unstructured.Unstructured)
					if !ok {
						return
					}
					fmt.Printf("deleted: %s\n", s.GetName())
				},
			})
		}

		informerSet.Factories = append(informerSet.Factories, factory)
	}

	return informerSet, nil
}

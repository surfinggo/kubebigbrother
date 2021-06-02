package informers

import (
	"fmt"
	"github.com/spongeprojects/kubebigbrother/pkg/utils/resourcebuilder"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/klog/v2"
	"time"
)

type Interface interface {
	// Start starts all informers registered,
	// Start is non-blocking, you should always call Shutdown before exit.
	Start(stopCh <-chan struct{}) error

	// Shutdown should be called before exit.
	Shutdown()
}

type InformerSet struct {
	ResourceBuilder   resourcebuilder.Interface
	Factories         []dynamicinformer.DynamicSharedInformerFactory
	ResourceInformers []*Informer
}

func (set *InformerSet) Start(stopCh <-chan struct{}) error {
	for i, factory := range set.Factories {
		klog.V(1).Infof("starting namespaced informer factory %d/%d: n%d",
			i+1, len(set.Factories), i)
		go factory.Start(stopCh)
	}

	klog.Info("waiting for caches to sync...")
	for i, factory := range set.Factories {
		for gvr, ok := range factory.WaitForCacheSync(stopCh) {
			if !ok {
				return fmt.Errorf(
					"timed out waiting for caches to sync, .Factories[%d][%s]",
					i, gvr)
			}
		}
	}
	klog.Info("caches synced, starting informers...")

	for i, resourceInformer := range set.ResourceInformers {
		klog.V(1).Infof("starting informer %d/%d: %s, workers: %d, resource: %s",
			i+1, len(set.ResourceInformers), resourceInformer.ID,
			resourceInformer.Workers, resourceInformer.Resource)
		for i := 0; i < resourceInformer.Workers; i++ {
			go wait.Until(resourceInformer.RunWorker, time.Second, stopCh)
		}
	}

	<-stopCh

	return nil
}

func (set *InformerSet) Shutdown() {
	for _, resourceInformer := range set.ResourceInformers {
		resourceInformer.ShutDown()
	}
}

// Setup setups new InformerSet
func Setup(config Config) (*InformerSet, error) {
	return NewBootstrapper(&config).buildInformerSet()
}

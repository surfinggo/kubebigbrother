package resourcebuilder

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"time"
)

// PersistentRESTClientGetter implements RESTClientGetter with persistent clients
type PersistentRESTClientGetter struct {
	RESTConfig               *rest.Config
	CachedDiscoveryInterface discovery.CachedDiscoveryInterface
	RESTMapper               meta.RESTMapper
}

// ToRESTConfig implements RESTClientGetter
func (g *PersistentRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return g.RESTConfig, nil
}

// ToDiscoveryClient implements RESTClientGetter
func (g *PersistentRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return g.CachedDiscoveryInterface, nil
}

// ToRESTMapper implements RESTClientGetter
func (g *PersistentRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return g.RESTMapper, nil
}

// NewPersistentRESTClientGetter creates new PersistentRESTClientGetter
func NewPersistentRESTClientGetter(kubeconfig string) (*PersistentRESTClientGetter, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "clientcmd.BuildConfigFromFlags error")
	}

	discoveryConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "clientcmd.BuildConfigFromFlags error")
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom resources) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	discoveryConfig.Burst = 100

	cacheDir, err := ioutil.TempDir("", "kubebigbrother")
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.TempDir error")
	}
	httpCacheDir := filepath.Join(cacheDir, "http")
	discoveryCacheDir := filepath.Join(cacheDir, "discovery")

	cachedDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(
		restConfig, discoveryCacheDir, httpCacheDir, 10*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "disk.NewCachedDiscoveryClientForConfig error")
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, cachedDiscoveryClient)

	return &PersistentRESTClientGetter{
		RESTConfig:               restConfig,
		CachedDiscoveryInterface: cachedDiscoveryClient,
		RESTMapper:               expander,
	}, nil
}

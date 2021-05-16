package informers

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
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

func (g *PersistentRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return g.RESTConfig, nil
}

func (g *PersistentRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return g.CachedDiscoveryInterface, nil
}

func (g *PersistentRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return g.RESTMapper, nil
}

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

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, cachedDiscoveryClient)

	return &PersistentRESTClientGetter{
		RESTConfig:               restConfig,
		CachedDiscoveryInterface: cachedDiscoveryClient,
		RESTMapper:               expander,
	}, nil
}

type ResourceBuilder struct {
	BaseBuilder *resource.Builder
}

// ParseGroupResource parses resource string as schema.GroupVersionResource,
func (b *ResourceBuilder) ParseGroupResource(resource string) (schema.GroupVersionResource, error) {
	r := b.BaseBuilder.Unstructured().SingleResourceType().
		ResourceTypeOrNameArgs(true, resource).Do()

	infos, err := r.Infos()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	if len(infos) != 1 {
		return schema.GroupVersionResource{}, errors.New("multiple info returned, expect 1")
	}
	return infos[0].Mapping.Resource, nil
}

func NewResourceBuilderFromClientGetter(clientGetter resource.RESTClientGetter) (*ResourceBuilder, error) {
	return &ResourceBuilder{
		BaseBuilder: resource.NewBuilder(clientGetter),
	}, nil
}

func NewResourceBuilder(kubeconfig string) (*ResourceBuilder, error) {
	clientGetter, err := NewPersistentRESTClientGetter(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "NewPersistentRESTClientGetter error")
	}
	return NewResourceBuilderFromClientGetter(clientGetter)
}

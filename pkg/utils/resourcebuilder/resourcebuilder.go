package resourcebuilder

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
)

// Interface is used to parse resource name
type Interface interface {
	ParseGroupResource(resourceArg string) (schema.GroupVersionResource, error)
}

// ResourceBuilder is used to parse resource name
type ResourceBuilder struct {
	ClientGetter resource.RESTClientGetter
}

// ParseGroupResource parses resource string as schema.GroupVersionResource,
func (b *ResourceBuilder) ParseGroupResource(resourceArg string) (schema.GroupVersionResource, error) {
	r := resource.NewBuilder(b.ClientGetter).Unstructured().SingleResourceType().
		ResourceTypeOrNameArgs(true, resourceArg).Do()

	infos, err := r.Infos()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	if len(infos) != 1 {
		return schema.GroupVersionResource{}, errors.New("multiple info returned, expect 1")
	}
	return infos[0].Mapping.Resource, nil
}

// NewFromClientGetter creates new ResourceBuilder
func NewFromClientGetter(clientGetter resource.RESTClientGetter) (*ResourceBuilder, error) {
	return &ResourceBuilder{
		ClientGetter: clientGetter,
	}, nil
}

// New creates new ResourceBuilder
func New(kubeconfig string) (*ResourceBuilder, error) {
	clientGetter, err := NewPersistentRESTClientGetter(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "NewPersistentRESTClientGetter error")
	}
	return NewFromClientGetter(clientGetter)
}

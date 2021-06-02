package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/magicconch"
)

type ResourceConfig struct {
	// Name is an unique value for this resource,
	// the controller, server and query command use Name to query events
	Name string `json:"name" yaml:"name"`

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
	UpdateOn []string `json:"updateOn,omitempty" yaml:"updateOn,omitempty"`

	// ChannelNames defines channels to send notification
	ChannelNames []ChannelName `json:"channelNames,omitempty" yaml:"channelNames,omitempty"`

	// ResyncPeriod is the resync period in reflectors for this resource
	ResyncPeriod string `json:"resyncPeriod,omitempty" yaml:"resyncPeriod,omitempty"`

	// Workers is the number of workers
	Workers int `json:"workers,omitempty" yaml:"workers,omitempty"`

	// MaxRetries is the max retry times
	MaxRetries int `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`
}

func (c *ResourceConfig) validate(
	namespaceIndex, index int, channelNames []string) error {
	if c.Name == "" {
		return errors.Errorf(
			"you must set a name for each resource to watch: .Namespaces[%d].Resources[%d]",
			namespaceIndex, index)
	}

	for _, name := range c.ChannelNames {
		if !magicconch.StringInSlice(string(name), channelNames) {
			return errors.Errorf(
				"non-exist channel name: %s in .Namespaces[%d].Resources[%d]",
				name, namespaceIndex, index)
		}
	}
	return nil
}

func (c *ResourceConfig) getChannelNames(
	namespaceDefault []ChannelName) []ChannelName {
	if len(c.ChannelNames) == 0 {
		return namespaceDefault
	}
	return c.ChannelNames
}

func (c *ResourceConfig) getWorkers(namespaceDefault int) int {
	if c.Workers < 1 {
		return namespaceDefault
	}
	return c.Workers
}

func (c *ResourceConfig) getMaxRetries(namespaceDefault int) int {
	if c.MaxRetries < 1 {
		return namespaceDefault
	}
	return c.MaxRetries
}

func (c *ResourceConfig) buildResyncPeriodFunc(
	namespaceDefault ResyncPeriodFunc) (ResyncPeriodFunc, error) {
	f, set, err := buildResyncPeriodFunc(c.ResyncPeriod)
	if err != nil {
		return nil, err
	}
	if !set {
		return namespaceDefault, nil
	}
	return f, nil
}

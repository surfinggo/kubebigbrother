package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/magicconch"
)

type NamespaceConfig struct {
	// Namespace is the namespace to watch, default to "", which means all namespaces
	Namespace string `json:"namespace" yaml:"namespace"`

	// Resources is the resources you want to watch
	Resources []ResourceConfig `json:"resources" yaml:"resources"`

	// DefaultChannelNames defines default channels in this namespace
	DefaultChannelNames []ChannelName `json:"defaultChannelNames,omitempty" yaml:"defaultChannelNames,omitempty"`

	// DefaultWorkers is the default number of workers in this namespace
	DefaultWorkers int `json:"defaultWorkers,omitempty" yaml:"defaultWorkers,omitempty"`

	// DefaultMaxRetries is the default max retry times in this namespace
	DefaultMaxRetries int `json:"defaultMaxRetries,omitempty" yaml:"defaultMaxRetries,omitempty"`

	// MinResyncPeriod is the resync period in reflectors in this namespace;
	// actual resync period will be random between MinResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod string `json:"minResyncPeriod,omitempty" yaml:"minResyncPeriod,omitempty"`
}

func (c *NamespaceConfig) validate(index int, channelNames []string) error {
	for _, name := range c.DefaultChannelNames {
		if !magicconch.StringInSlice(string(name), channelNames) {
			return errors.Errorf(
				"non-exist channel name: %s in .Namespaces[%d]", name, index)
		}
	}
	for i, resource := range c.Resources {
		if err := resource.validate(index, i, channelNames); err != nil {
			return err
		}
	}
	return nil
}

func (c *NamespaceConfig) getDefaultChannelNames(
	globalDefault []ChannelName) []ChannelName {
	if len(c.DefaultChannelNames) == 0 {
		return globalDefault
	}
	return c.DefaultChannelNames
}

func (c *NamespaceConfig) getDefaultWorkers(globalDefault int) int {
	if c.DefaultWorkers < 1 {
		return globalDefault
	}
	return c.DefaultWorkers
}

func (c *NamespaceConfig) getDefaultMaxRetries(globalDefault int) int {
	if c.DefaultMaxRetries < 1 {
		return globalDefault
	}
	return c.DefaultMaxRetries
}

func (c *NamespaceConfig) buildResyncPeriodFunc(
	defaultFunc ResyncPeriodFunc) (ResyncPeriodFunc, error) {
	f, set, err := buildResyncPeriodFunc(c.MinResyncPeriod)
	if err != nil {
		return nil, err
	}
	if !set {
		return defaultFunc, nil
	}
	return f, nil
}

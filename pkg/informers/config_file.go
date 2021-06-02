package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/magicconch"
)

// ConfigFile is struct of informers config file,
// order of fields affects order of fields in JSON
type ConfigFile struct {
	// Channels defines channels that receive notifications
	Channels map[ChannelName]ChannelConfig `json:"channels" yaml:"channels"`

	// Namespaces defines namespaces and resources to watch
	Namespaces []NamespaceConfig `json:"namespaces" yaml:"namespaces"`

	// DefaultChannelNames defines default channels
	DefaultChannelNames []ChannelName `json:"defaultChannelNames,omitempty" yaml:"defaultChannelNames,omitempty"`

	// DefaultWorkers is the default number of workers
	DefaultWorkers int `json:"defaultWorkers,omitempty" yaml:"defaultWorkers,omitempty"`

	// DefaultMaxRetries is the default max retry times
	DefaultMaxRetries int `json:"defaultMaxRetries,omitempty" yaml:"defaultMaxRetries,omitempty"`

	// MinResyncPeriod is the resync period in reflectors;
	// actual resync period will be random between MinResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod string `json:"minResyncPeriod,omitempty" yaml:"minResyncPeriod,omitempty"`
}

// Validate validates ConfigFile values
func (c *ConfigFile) Validate() error {
	var channelNames []string
	for name, channel := range c.Channels {
		channelNames = append(channelNames, string(name))

		switch channel.Type {
		case channels.ChannelTypeCallback:
			if channel.Callback == nil {
				return errors.Errorf(
					"config missing for callback channel, name: %s", name)
			}
		case channels.ChannelTypeDingtalk:
			if channel.Dingtalk == nil {
				return errors.Errorf(
					"config missing for Dingtalk channel, name: %s", name)
			}
		case channels.ChannelTypeFlock:
			if channel.Flock == nil {
				return errors.Errorf(
					"config missing for Flock channel, name: %s", name)
			}
		case channels.ChannelTypePrint:
			if channel.Print == nil {
				return errors.Errorf(
					"config missing for print channel, name: %s", name)
			}
		case channels.ChannelTypeSlack:
			if channel.Slack == nil {
				return errors.Errorf(
					"config missing for Slack channel, name: %s", name)
			}
		case channels.ChannelTypeTelegram:
			if channel.Telegram == nil {
				return errors.Errorf(
					"config missing for Telegram channel, name: %s", name)
			}
		}
	}
	for _, name := range c.DefaultChannelNames {
		if !magicconch.StringInSlice(string(name), channelNames) {
			return errors.Errorf(
				"non-exist channel name: %s in default namespaces", name)
		}
	}
	for i, namespace := range c.Namespaces {
		if err := namespace.validate(i, channelNames); err != nil {
			return err
		}
	}
	return nil
}

// getDefaultWorkers returns global default workers with default value
func (c *ConfigFile) getDefaultWorkers() int {
	if c.DefaultWorkers < 1 {
		return 3
	}
	return c.DefaultWorkers
}

// getDefaultMaxRetries returns global default max retries with default value
func (c *ConfigFile) getDefaultMaxRetries() int {
	if c.DefaultMaxRetries < 1 {
		return 3
	}
	return c.DefaultMaxRetries
}

func (c *ConfigFile) buildResyncPeriodFunc() (f ResyncPeriodFunc, err error) {
	if c.MinResyncPeriod == "" {
		c.MinResyncPeriod = "12h"
	}
	f, _, err = buildResyncPeriodFunc(c.MinResyncPeriod)
	return f, err
}

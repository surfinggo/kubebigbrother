package informers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// ChannelGroupConfig is config for ChannelGroup, read from config file
type ChannelGroupConfig []channels.ChannelName

// ChannelTelegramConfig is config for ChannelTelegram, read from config file
type ChannelTelegramConfig struct {
	Token           string  `json:"token" yaml:"token"`
	Recipients      []int64 `json:"recipients" yaml:"recipients"`
	Proxy           string  `json:"proxy" yaml:"proxy"`
	AddedTemplate   string  `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string  `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string  `json:"updatedTemplate" yaml:"updatedTemplate"`
}

// ChannelCallbackConfig is config for ChannelCallback, read from config file
type ChannelCallbackConfig struct {
	URL string `json:"url" yaml:"url"`
}

// ChannelPrintConfig is config for ChannelPrint, read from config file
type ChannelPrintConfig struct {
	Writer          string `json:"writer" yaml:"writer"`
	AddedTemplate   string `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string `json:"updatedTemplate" yaml:"updatedTemplate"`
}

// ChannelConfig defines a channel to receive notifications
type ChannelConfig struct {
	// Type is the type of the channel
	Type channels.ChannelType `json:"type" yaml:"type"`

	Callback *ChannelCallbackConfig `json:"callback" yaml:"callback"`
	Group    *ChannelGroupConfig    `json:"group" yaml:"group"`
	Print    *ChannelPrintConfig    `json:"print" yaml:"print"`
	Telegram *ChannelTelegramConfig `json:"telegram" yaml:"telegram"`
}

type ResourceConfig struct {
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

	// ChannelNames defines channels to send notification
	ChannelNames []channels.ChannelName `json:"channelNames" yaml:"channelNames"`

	// ResyncPeriod is the resync period in reflectors for this resource
	ResyncPeriod string `json:"resyncPeriod" yaml:"resyncPeriod"`

	// Workers is the number of workers
	Workers int `json:"workers" yaml:"workers"`
}

func (c *ResourceConfig) buildResyncPeriodFuncWithDefault(defaultFunc resyncPeriodFunc) (resyncPeriodFunc, error) {
	f, set, err := c.buildResyncPeriodFunc()
	if err != nil {
		return nil, err
	}
	if !set {
		return defaultFunc, nil
	}
	return f, nil
}

func (c *ResourceConfig) buildResyncPeriodFunc() (f func() time.Duration, set bool, err error) {
	return buildResyncPeriodFunc(c.ResyncPeriod)
}

type NamespaceConfig struct {
	// Namespace is the namespace to watch, default to "", which means all namespaces
	Namespace string `json:"namespace" yaml:"namespace"`

	// Resources is the resources you want to watch
	Resources []ResourceConfig `json:"resources" yaml:"resources"`

	// DefaultChannelNames defines default channels in this namespace
	DefaultChannelNames []channels.ChannelName `json:"defaultChannelNames" yaml:"defaultChannelNames"`

	// DefaultWorkers is the default number of workers in this namespace
	DefaultWorkers int `json:"defaultWorkers" yaml:"defaultWorkers"`

	// MinResyncPeriod is the resync period in reflectors in this namespace;
	// actual resync period will be random between MinResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod string `json:"minResyncPeriod" yaml:"minResyncPeriod"`
}

func (c *NamespaceConfig) buildResyncPeriodFuncWithDefault(defaultFunc resyncPeriodFunc) (resyncPeriodFunc, error) {
	f, set, err := c.buildResyncPeriodFunc()
	if err != nil {
		return nil, err
	}
	if !set {
		return defaultFunc, nil
	}
	return f, nil
}

func (c *NamespaceConfig) buildResyncPeriodFunc() (f resyncPeriodFunc, set bool, err error) {
	return buildResyncPeriodFunc(c.MinResyncPeriod)
}

type Config struct {
	// Namespaces defines namespaces and resources to watch
	Namespaces []NamespaceConfig `json:"namespaces" yaml:"namespaces"`

	// Channels defines channels that receive notifications
	Channels map[channels.ChannelName]ChannelConfig `json:"channels" yaml:"channels"`

	// DefaultChannelNames defines default channels
	DefaultChannelNames []channels.ChannelName `json:"defaultChannelNames" yaml:"defaultChannelNames"`

	// DefaultWorkers is the default number of workers
	DefaultWorkers int `json:"defaultWorkers" yaml:"defaultWorkers"`

	// MinResyncPeriod is the resync period in reflectors;
	// actual resync period will be random between MinResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod string `json:"minResyncPeriod" yaml:"minResyncPeriod"`
}

func (c *Config) buildResyncPeriodFunc() (f resyncPeriodFunc, err error) {
	if c.MinResyncPeriod == "" {
		c.MinResyncPeriod = "12h"
	}
	f, _, err = buildResyncPeriodFunc(c.MinResyncPeriod)
	return f, err
}

type resyncPeriodFunc func() time.Duration

func buildResyncPeriodFunc(resyncPeriod string) (f resyncPeriodFunc, set bool, err error) {
	duration, set, err := parseResyncPeriod(resyncPeriod)
	if err != nil {
		return nil, false, err
	}
	if !set {
		return nil, false, nil
	}
	durationFloat := float64(duration.Nanoseconds())
	// generate time.Duration between duration and 2*duration
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(durationFloat * factor)
	}, true, nil
}

func parseResyncPeriod(resyncPeriod string) (f time.Duration, set bool, err error) {
	if resyncPeriod == "" {
		return 0, false, nil
	}
	duration, err := time.ParseDuration(resyncPeriod)
	if err != nil {
		return 0, false, errors.Wrap(err, "time.ParseDuration error")
	}
	return duration, true, nil
}

// LoadConfigFromFile loads config from file
func LoadConfigFromFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open error")
	}
	var config Config
	switch t := strings.ToLower(path.Ext(file)); t {
	case ".json":
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			return nil, errors.Wrap(err, "json decode error")
		}
	case ".yaml":
		err = yaml.NewDecoder(f).Decode(&config)
		if err != nil {
			return nil, errors.Wrap(err, "yaml decode error")
		}
	default:
		return nil, errors.Errorf("unsupported file type: %s", t)
	}
	return &config, nil
}

func buildChannelFromConfig(config *ChannelConfig) (channels.Channel, error) {
	switch config.Type {
	case channels.ChannelTypeCallback:
		return channels.NewChannelCallback(config.Callback.URL)
	case channels.ChannelTypeGroup:
		return channels.NewChannelGroup()
	case channels.ChannelTypePrint:
		return channels.NewChannelPrint(config.Print.Writer,
			config.Print.AddedTemplate,
			config.Print.DeletedTemplate,
			config.Print.UpdatedTemplate)
	case channels.ChannelTypeTelegram:
		return channels.NewChannelTelegram(&channels.ChannelTelegramConfig{
			Token:           config.Telegram.Token,
			Recipients:      config.Telegram.Recipients,
			Proxy:           config.Telegram.Proxy,
			AddedTemplate:   config.Telegram.AddedTemplate,
			DeletedTemplate: config.Telegram.DeletedTemplate,
			UpdatedTemplate: config.Telegram.UpdatedTemplate,
		})
	default:
		return nil, errors.Errorf("unsupported channel type: %s", config.Type)
	}
}

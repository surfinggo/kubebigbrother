package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
)

// ChannelName is an alias to channels.ChannelName
type ChannelName = channels.ChannelName

// ChannelCallbackConfig is config for ChannelCallback, read from config file
type ChannelCallbackConfig struct {
	Method          string `json:"method,omitempty" yaml:"method,omitempty"`
	URL             string `json:"url" yaml:"url"`
	Proxy           string `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	UseTemplate     bool   `json:"useTemplate,omitempty" yaml:"useTemplate,omitempty"`
	AddedTemplate   string `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelCallbackConfig) setupChannel() (*channels.ChannelCallback, error) {
	return channels.NewChannelCallback(&channels.ChannelCallbackConfig{
		Method:          c.Method,
		Proxy:           c.Proxy,
		URL:             c.URL,
		UseTemplate:     c.UseTemplate,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelDingtalkConfig is config for ChannelDingtalk, read from config file
type ChannelDingtalkConfig struct {
	WebhookURL      string   `json:"webhookURL" yaml:"webhookURL"`
	Proxy           string   `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	AtMobiles       []string `json:"atMobiles,omitempty" yaml:"atMobiles,omitempty"`
	AtAll           bool     `json:"atAll,omitempty" yaml:"atAll,omitempty"`
	AddedTemplate   string   `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string   `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string   `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelDingtalkConfig) setupChannel() (*channels.ChannelDingtalk, error) {
	return channels.NewChannelDingtalk(&channels.ChannelDingtalkConfig{
		WebhookURL:      c.WebhookURL,
		Proxy:           c.Proxy,
		AtMobiles:       c.AtMobiles,
		AtAll:           c.AtAll,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelFlockConfig is config for ChannelFlock, read from config file
type ChannelFlockConfig struct {
	URL             string `json:"url" yaml:"url"`
	Proxy           string `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	TitleTemplate   string `json:"titleTemplate,omitempty" yaml:"titleTemplate,omitempty"`
	AddedTemplate   string `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelFlockConfig) setupChannel() (*channels.ChannelFlock, error) {
	return channels.NewChannelFlock(&channels.ChannelFlockConfig{
		URL:             c.URL,
		Proxy:           c.Proxy,
		TitleTemplate:   c.TitleTemplate,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelPrintConfig is config for ChannelPrint, read from config file
type ChannelPrintConfig struct {
	Writer          string `json:"writer" yaml:"writer"`
	UseTemplate     bool   `json:"useTemplate,omitempty" yaml:"useTemplate,omitempty"`
	AddedTemplate   string `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelPrintConfig) setupChannel() (*channels.ChannelPrint, error) {
	return channels.NewChannelPrint(&channels.ChannelPrintConfig{
		Writer:          c.Writer,
		UseTemplate:     c.UseTemplate,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelSlackConfig is config for ChannelSlack, read from config file
type ChannelSlackConfig struct {
	Token           string `json:"token" yaml:"token"`
	Proxy           string `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	WebhookURL      string `json:"webhookURL" yaml:"webhookURL"`
	TitleTemplate   string `json:"titleTemplate,omitempty" yaml:"titleTemplate,omitempty"`
	AddedTemplate   string `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelSlackConfig) setupChannel() (*channels.ChannelSlack, error) {
	return channels.NewChannelSlack(&channels.ChannelSlackConfig{
		Token:           c.Token,
		Proxy:           c.Proxy,
		WebhookURL:      c.WebhookURL,
		TitleTemplate:   c.TitleTemplate,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelTelegramConfig is config for ChannelTelegram, read from config file
type ChannelTelegramConfig struct {
	Token           string   `json:"token" yaml:"token"`
	ChatIDs         []string `json:"chatIDs" yaml:"chatIDs"`
	Proxy           string   `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	AddedTemplate   string   `json:"addedTemplate,omitempty" yaml:"addedTemplate,omitempty"`
	DeletedTemplate string   `json:"deletedTemplate,omitempty" yaml:"deletedTemplate,omitempty"`
	UpdatedTemplate string   `json:"updatedTemplate,omitempty" yaml:"updatedTemplate,omitempty"`
}

func (c *ChannelTelegramConfig) setupChannel() (*channels.ChannelTelegram, error) {
	return channels.NewChannelTelegram(&channels.ChannelTelegramConfig{
		Token:           c.Token,
		ChatIDs:         c.ChatIDs,
		Proxy:           c.Proxy,
		AddedTemplate:   c.AddedTemplate,
		DeletedTemplate: c.DeletedTemplate,
		UpdatedTemplate: c.UpdatedTemplate,
	})
}

// ChannelConfig defines a channel to receive notifications
type ChannelConfig struct {
	// Type is the type of the channel
	Type channels.ChannelType `json:"type" yaml:"type"`

	Callback *ChannelCallbackConfig `json:"callback,omitempty" yaml:"callback,omitempty"`
	Dingtalk *ChannelDingtalkConfig `json:"dingtalk,omitempty" yaml:"dingtalk,omitempty"`
	Flock    *ChannelFlockConfig    `json:"flock,omitempty" yaml:"flock,omitempty"`
	Print    *ChannelPrintConfig    `json:"print,omitempty" yaml:"print,omitempty"`
	Slack    *ChannelSlackConfig    `json:"slack,omitempty" yaml:"slack,omitempty"`
	Telegram *ChannelTelegramConfig `json:"telegram,omitempty" yaml:"telegram,omitempty"`
}

func (c *ChannelConfig) setupChannel() (channels.Channel, error) {
	switch c.Type {
	case channels.ChannelTypeCallback:
		return c.Callback.setupChannel()
	case channels.ChannelTypeDingtalk:
		return c.Dingtalk.setupChannel()
	case channels.ChannelTypeFlock:
		return c.Flock.setupChannel()
	case channels.ChannelTypePrint:
		return c.Print.setupChannel()
	case channels.ChannelTypeSlack:
		return c.Slack.setupChannel()
	case channels.ChannelTypeTelegram:
		return c.Telegram.setupChannel()
	default:
		return nil, errors.Errorf("unsupported channel type: %s", c.Type)
	}
}

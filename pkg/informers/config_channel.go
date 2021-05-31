package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
)

// ChannelName is an alias to channels.ChannelName
type ChannelName = channels.ChannelName

// ChannelCallbackConfig is config for ChannelCallback, read from config file
type ChannelCallbackConfig struct {
	Method          string `json:"method" yaml:"method"`
	URL             string `json:"url" yaml:"url"`
	Proxy           string `json:"proxy" yaml:"proxy"`
	UseTemplate     bool   `json:"useTemplate" yaml:"useTemplate"`
	AddedTemplate   string `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string `json:"updatedTemplate" yaml:"updatedTemplate"`
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

// ChannelFlockConfig is config for ChannelFlock, read from config file
type ChannelFlockConfig struct {
	URL             string `json:"url" yaml:"url"`
	Proxy           string `json:"proxy" yaml:"proxy"`
	TitleTemplate   string `json:"titleTemplate" yaml:"titleTemplate"`
	AddedTemplate   string `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string `json:"updatedTemplate" yaml:"updatedTemplate"`
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
	UseTemplate     bool   `json:"useTemplate" yaml:"useTemplate"`
	AddedTemplate   string `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string `json:"updatedTemplate" yaml:"updatedTemplate"`
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
	Proxy           string `json:"proxy" yaml:"proxy"`
	WebhookURL      string `json:"webhookURL" yaml:"webhookURL"`
	TitleTemplate   string `json:"titleTemplate" yaml:"titleTemplate"`
	AddedTemplate   string `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string `json:"updatedTemplate" yaml:"updatedTemplate"`
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
	Proxy           string   `json:"proxy" yaml:"proxy"`
	AddedTemplate   string   `json:"addedTemplate" yaml:"addedTemplate"`
	DeletedTemplate string   `json:"deletedTemplate" yaml:"deletedTemplate"`
	UpdatedTemplate string   `json:"updatedTemplate" yaml:"updatedTemplate"`
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

	Callback *ChannelCallbackConfig `json:"callback" yaml:"callback"`
	Flock    *ChannelFlockConfig    `json:"flock" yaml:"flock"`
	Print    *ChannelPrintConfig    `json:"print" yaml:"print"`
	Slack    *ChannelSlackConfig    `json:"slack" yaml:"slack"`
	Telegram *ChannelTelegramConfig `json:"telegram" yaml:"telegram"`
}

func (c *ChannelConfig) setupChannel() (channels.Channel, error) {
	switch c.Type {
	case channels.ChannelTypeCallback:
		return c.Callback.setupChannel()
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

package informers

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"k8s.io/klog/v2"
	"net/http"
)

// ChannelName is name of channel
type ChannelName string

// ChannelType is type of channel
type ChannelType string

const (
	ChannelTypeGroup    = "group"
	ChannelTypeTelegram = "telegram"
	ChannelTypeCallback = "callback"
	ChannelTypePrint    = "print"
)

// ChannelGroupConfig is a special type of channel, which includes a slice of channels
type ChannelGroupConfig []ChannelName

// ChannelTelegramConfig is the Telegram channel
type ChannelTelegramConfig struct {
	Token string `json:"token" yaml:"token"`
}

// ChannelCallbackConfig is the callback channel
type ChannelCallbackConfig struct {
	URL string `json:"url" yaml:"url"`
}

// ChannelPrintConfig is the channel to print event to writer
type ChannelPrintConfig struct {
	Out io.Writer
}

type ChannelMap map[ChannelName]Channel

type Channel interface {
	Handle(e *Event) error
}

func BuildChannelFromConfig(config *ChannelConfig) (Channel, error) {
	switch config.Type {
	case ChannelTypeCallback:
		return &ChannelCallback{
			Config: config.Callback,
			Client: http.DefaultClient,
		}, nil
	case ChannelTypeGroup:
		return &ChannelGroup{
			Config:   config.Group,
			Channels: nil, // TODO: set channels
		}, nil
	case ChannelTypePrint:
		return &ChannelPrint{
			Config: config.Print,
		}, nil
	case ChannelTypeTelegram:
		return &ChannelTelegram{
			Config: config.Telegram,
		}, nil
	default:
		return nil, errors.Errorf("unsupported channel type: %s", config.Type)
	}
}

// ChannelCallback is the callback channel
type ChannelCallback struct {
	Config *ChannelCallbackConfig

	Client *http.Client
}

func (c *ChannelCallback) Handle(e *Event) error {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(e)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	resp, err := c.Client.Post(c.Config.URL, "application/json", body)
	if err != nil {
		return errors.Wrap(err, "send request error")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("non-200 code returned: %d", resp.StatusCode)
	}
	return nil
}

// ChannelGroup is a special type of channel, which includes a slice of channels
type ChannelGroup struct {
	Config *ChannelGroupConfig

	Channels []Channel
}

func (c *ChannelGroup) Handle(e *Event) error {
	for _, channel := range c.Channels {
		err := channel.Handle(e)
		if err != nil {
			klog.V(6).Info(errors.Wrap(err, "handle event error"))
		}
	}
	// TODO: if an error is returned, all channels in c.ChannelNames will be retried.
	// need to find a way to handle error but only retry failed channels.
	return nil
}

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Config *ChannelPrintConfig
}

func (c *ChannelPrint) Handle(e *Event) error {
	err := json.NewEncoder(c.Config.Out).Encode(e)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	return nil
}

// ChannelTelegram is the Telegram channel
type ChannelTelegram struct {
	Config *ChannelTelegramConfig
}

func (c *ChannelTelegram) Handle(_ *Event) error {
	panic("wip")
}

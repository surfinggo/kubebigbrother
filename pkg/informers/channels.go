package informers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"os"
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

// ChannelGroupConfig is config for ChannelGroup, read from config file
type ChannelGroupConfig []ChannelName

// ChannelTelegramConfig is config for ChannelTelegram, read from config file
type ChannelTelegramConfig struct {
	Token string `json:"token" yaml:"token"`
}

// ChannelCallbackConfig is config for ChannelCallback, read from config file
type ChannelCallbackConfig struct {
	URL string `json:"url" yaml:"url"`
}

const (
	PrintWriterStdout = "stdout"
)

// ChannelPrintConfig is config for ChannelPrint, read from config file
type ChannelPrintConfig struct {
	Writer string
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
		var writer io.Writer
		switch config.Print.Writer {
		case PrintWriterStdout, "":
			writer = os.Stdout
		default:
			return nil, errors.Errorf("unsupported writer: %s", config.Print.Writer)
		}
		return &ChannelPrint{
			Config: config.Print,
			Writer: writer,
			// TODO: make WriteFunc configurable
			WriteFunc: func(e *Event, w io.Writer) error {
				t := fmt.Sprintf("[%s] %s\n", e.Type, NamespaceKey(e.Obj))
				_, err := w.Write([]byte(t))
				return err
			},
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
			klog.Warning(errors.Wrap(err, "handle event error"))
		}
	}
	// TODO: if an error is returned, all channels in c.ChannelNames will be retried.
	// need to find a way to handle error but only retry failed channels.
	return nil
}

// ChannelPrint is the channel to print event to writer
type ChannelPrint struct {
	Config *ChannelPrintConfig

	Writer    io.Writer
	WriteFunc func(*Event, io.Writer) error
}

func (c *ChannelPrint) Handle(e *Event) error {
	if c.WriteFunc != nil {
		return c.WriteFunc(e, c.Writer)
	}
	err := json.NewEncoder(c.Writer).Encode(e)
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

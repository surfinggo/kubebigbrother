package channels

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
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

type ChannelMap map[ChannelName]Channel

type Channel interface {
	Handle(e *event.Event) error
}

// ChannelCallback is the callback channel
type ChannelCallback struct {
	Client *http.Client
	URL    string
}

func (c *ChannelCallback) Handle(e *event.Event) error {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(e)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	resp, err := c.Client.Post(c.URL, "application/json", body)
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
	Channels []Channel
}

func (c *ChannelGroup) Handle(e *event.Event) error {
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
	Writer    io.Writer
	WriteFunc func(*event.Event, io.Writer) error
}

func (c *ChannelPrint) Handle(e *event.Event) error {
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
}

func (c *ChannelTelegram) Handle(_ *event.Event) error {
	panic("wip")
}

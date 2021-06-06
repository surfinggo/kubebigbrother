package informers

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/stores/event_store"
	"time"
)

type Config struct {
	Kubeconfig          string
	DefaultWorkers      int
	DefaultMaxRetries   int
	DefaultChannelNames []string
	MinResyncPeriod     time.Duration
	JustWatch           bool
	EventStore          event_store.Interface
}

func (c *Config) Validate() error {
	if !c.JustWatch && c.EventStore == nil {
		return errors.New("event store cannot be nil when not just watching")
	}
	return nil
}

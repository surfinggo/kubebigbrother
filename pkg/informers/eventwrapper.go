package informers

import (
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	"github.com/spongeprojects/kubebigbrother/pkg/event"
)

// ChannelToProcess defines a channel to process
type ChannelToProcess struct {
	ChannelName         string
	EventProcessContext *channels.EventProcessContext
}

// eventWrapper wraps an event to process,
// ChannelsToProcess should be set only once when creating eventWrapper,
// it is updated as the process go on.
type eventWrapper struct {
	*event.Event

	// ChannelsToProcess defines channels to process
	ChannelsToProcess []ChannelToProcess
}

func (s *InformerSet) wrap(e *event.Event, channelNames []string) *eventWrapper {
	if s.JustWatch {
		channelNames = []string{channelNamePrintToStdout}
	}

	var channelsToProcess []ChannelToProcess
	for _, name := range channelNames {
		if channel, ok := s.ChannelMap[name]; ok {
			channelsToProcess = append(channelsToProcess, ChannelToProcess{
				ChannelName:         name,
				EventProcessContext: channel.NewEventProcessContext(e),
			})
		}
	}

	return &eventWrapper{
		Event:             e,
		ChannelsToProcess: channelsToProcess,
	}
}

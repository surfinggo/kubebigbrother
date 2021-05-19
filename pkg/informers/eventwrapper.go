package informers

type EventWrapper struct {
	Event *Event

	// ChannelNames is channels to process,
	// name will be removed from slice after processed successfully
	ChannelNames []ChannelName
}

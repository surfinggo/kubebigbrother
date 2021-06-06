package informers

import (
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/channels"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
)

func (s *InformerSet) RunChannelWorker() {
	for s.processNextChannel() {
	}
}

// processNextChannel waits and processes items in the queue
func (s *InformerSet) processNextChannel() bool {
	// block until an key arrives or queue shutdown
	obj, shutdown := s.ChannelQueue.Get()
	if shutdown {
		return false
	}
	key := obj.(string)

	if klog.V(5).Enabled() {
		klog.Infof("[channel] [%s try] key pop from queue: [%s]",
			humanize.Ordinal(s.ChannelQueue.NumRequeues(key)+1), key)
	}

	// we need to mark key as completed whether success or fail
	defer s.ChannelQueue.Done(key)

	result := s.processChannel(key)
	s.handleChannelErr(key, result)

	return true
}

// processChannel process an item synchronously
func (s *InformerSet) processChannel(key string) error {
	channel, err := s.ChannelLister.Get(key)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Infof("channel deleted: %s", key)
			delete(s.ChannelMap, key)
			return nil
		}
		return errors.Wrap(err, "get channel error")
	}

	var channelInstance channels.Channel
	switch channel.Spec.Type {
	case channels.ChannelTypeCallback:
		if channel.Spec.Callback == nil {
			return errors.Errorf("config missing for callback channel")
		}
		channelInstance, err = channels.NewChannelCallback(channel.Spec.Callback)
	case channels.ChannelTypeDingtalk:
		if channel.Spec.Dingtalk == nil {
			return errors.Errorf("config missing for Dingtalk channel")
		}
		channelInstance, err = channels.NewChannelDingtalk(channel.Spec.Dingtalk)
	case channels.ChannelTypeFlock:
		if channel.Spec.Flock == nil {
			return errors.Errorf("config missing for Flock channel")
		}
		channelInstance, err = channels.NewChannelFlock(channel.Spec.Flock)
	case channels.ChannelTypePrint:
		if channel.Spec.Print == nil {
			return errors.Errorf("config missing for print channel")
		}
		channelInstance, err = channels.NewChannelPrint(channel.Spec.Print)
	case channels.ChannelTypeSlack:
		if channel.Spec.Slack == nil {
			return errors.Errorf("config missing for Slack channel")
		}
		channelInstance, err = channels.NewChannelSlack(channel.Spec.Slack)
	case channels.ChannelTypeTelegram:
		if channel.Spec.Telegram == nil {
			return errors.Errorf("config missing for Telegram channel")
		}
		channelInstance, err = channels.NewChannelTelegram(channel.Spec.Telegram)
	default:
		return errors.Errorf("unsupported channel type: %s", channel.Spec.Type)
	}
	if err != nil {
		return errors.Wrap(err, "create channel instance error")
	}

	klog.V(2).Infof("[channel] channel added/updated: %s", key)
	s.ChannelMap[key] = channelInstance
	return nil
}

// handleChannelErr checks the result, schedules retry if needed
func (s *InformerSet) handleChannelErr(key string, result error) {
	if result == nil {
		if klog.V(2).Enabled() {
			klog.Infof("[channel] [%s try] key processed: [%s]",
				humanize.Ordinal(s.ChannelQueue.NumRequeues(key)+1), key)
		}
		// clear retry counter after success
		s.ChannelQueue.Forget(key)
		return
	}

	// retrying
	klog.Warningf("[channel] [%s try] error processing: [%s]: %s, will be retried",
		humanize.Ordinal(s.ChannelQueue.NumRequeues(key)+1), key, result)
	s.ChannelQueue.AddRateLimited(key)
}

package telegram

import (
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"k8s.io/klog/v2"
	"time"
)

func NewBot(token string) (*tb.Bot, error) {
	var bot *tb.Bot
	var err error
	count := 1
	for {
		klog.Infof("[%d times] trying to create telegram bot...", count)
		bot, err = tb.NewBot(tb.Settings{
			Token: token,
		})
		if err == nil {
			klog.Infof("[%d times] telegram bot created", count)
			break
		}
		if count >= 10 {
			return nil, errors.Wrapf(err, "[%d times] create telegram bot error, max retry exceeded", count)
		}
		klog.Warning(errors.Wrapf(err, "[%d times] create telegram bot error, retrying...", count))
		time.Sleep(2 * time.Second)
		count++
	}
	return bot, nil
}

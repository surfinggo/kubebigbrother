package telegram

import (
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"time"
)

// NewBot creates a new Telegram bot
func NewBot(token, proxy string) (*tb.Bot, error) {
	var httpClient *http.Client
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", proxy)
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else {
		httpClient = http.DefaultClient
	}
	setting := &tb.Settings{
		Token:  token,
		Client: httpClient,
	}

	count := 1
	for {
		t := humanize.Ordinal(count)

		klog.V(1).Infof("[%s time] trying to connect Telegram...", t)

		bot, err := tb.NewBot(*setting)
		if err == nil {
			klog.V(1).Infof("[%s time] Telegram connected", t)
			return bot, nil
		}

		if count >= 10 {
			return nil, errors.Wrapf(err,
				"[%s time] connect Telegram error, max retry exceeded", t)
		}

		klog.V(1).Info(errors.Wrapf(err,
			"[%s time] connect Telegram error, retrying...", t))

		time.Sleep(2 * time.Second)

		count++
	}
}

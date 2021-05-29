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
	if len(token) < 40 {
		return nil, errors.New("invalid token, too short")
	}

	klog.V(1).Infof("using Telegram token: %s...", token[:15])

	var httpClient *http.Client
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid proxy url: %s", proxy)
		}

		klog.V(1).Infof("connect to Telegram via proxy: %s", proxyUrl)

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

		klog.V(1).Infof("[%s time] trying to connect to Telegram...", t)

		bot, err := tb.NewBot(*setting)
		if err == nil {
			klog.V(1).Infof("[%s time] Telegram connected", t)
			return bot, nil
		}

		klog.V(1).Infof("[%s time] connect to Telegram error: %s", t, err)

		if count >= 10 {
			return nil, errors.Wrap(err,
				"connect to Telegram error, max retry exceeded")
		}

		time.Sleep(2 * time.Second)

		count++
	}
}

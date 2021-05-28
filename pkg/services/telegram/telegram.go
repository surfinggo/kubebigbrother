package telegram

import (
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
		t := "time"
		if count > 1 {
			t = "times"
		}

		klog.V(1).Infof("[%d %s] trying to connect Telegram...", count, t)

		bot, err := tb.NewBot(*setting)
		if err == nil {
			klog.V(1).Infof("[%d %s] Telegram connected", count, t)
			return bot, nil
		}

		if count >= 10 {
			return nil, errors.Wrapf(err,
				"[%d %s] connect Telegram error, max retry exceeded",
				count, t)
		}

		klog.V(1).Info(errors.Wrapf(err,
			"[%d %s] connect Telegram error, retrying...",
			count, t))

		time.Sleep(2 * time.Second)

		count++
	}
}

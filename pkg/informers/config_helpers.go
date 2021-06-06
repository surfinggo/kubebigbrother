package informers

import (
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

// ResyncPeriodFunc is a function to build resync period (time.Duration)
type ResyncPeriodFunc func() time.Duration

func buildResyncPeriodFuncByDuration(resyncPeriod time.Duration) (f ResyncPeriodFunc) {
	durationFloat := float64(resyncPeriod)
	// generate time.Duration between duration and 2*duration
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(durationFloat * factor)
	}
}

func parseResyncPeriod(resyncPeriod string) (d time.Duration, set bool, err error) {
	if resyncPeriod == "" {
		return 0, false, nil
	}
	duration, err := time.ParseDuration(resyncPeriod)
	if err != nil {
		return 0, false, errors.Wrap(err, "time.ParseDuration error")
	}
	return duration, true, nil
}

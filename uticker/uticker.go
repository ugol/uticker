package uticker

import (
	"time"
)

type TickerConfig struct {
	Duration       time.Duration
	ImmediateStart bool
}

type UTicker struct {
	C      chan time.Time
	ticker time.Ticker
	config TickerConfig
}

func WithImmediateStart(c *TickerConfig) *TickerConfig {
	c.ImmediateStart = true
	return c
}

func WithDuration(c *TickerConfig) *TickerConfig {
	d := 1 * time.Second
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	c.Duration = d
	return c
}

func NewUTicker(options ...func(*TickerConfig) *TickerConfig) *UTicker {

	t := &UTicker{
		C: make(chan time.Time, 1),
		config: TickerConfig{
			Duration:       1 * time.Second,
			ImmediateStart: false,
		},
	}

	if options != nil {
		for _, option := range options {
			option(&t.config)
		}
	}

	t.ticker = *time.NewTicker(t.config.Duration)

	go func() {
		if t.config.ImmediateStart {
			t.C <- time.Now()
		}
		for {
			select {
			case <-t.ticker.C:
				t.C <- time.Now()
			}
		}
	}()
	return t
}

func (t *UTicker) Stop() {

	//t.ticker.Stop()
	//close(t.C)

}

func (t *UTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	t.ticker.Reset(d)
}

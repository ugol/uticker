package uticker

import (
	"time"
)

type UTicker struct {
	C              chan time.Time
	ticker         *time.Ticker
	Duration       time.Duration
	ImmediateStart bool
}

func WithImmediateStart() func(*UTicker) {
	return func(t *UTicker) {
		t.ImmediateStart = true
	}
}

func WithDuration(d time.Duration) func(*UTicker) {
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	return func(t *UTicker) {
		t.Duration = d
	}
}

func NewUTicker(options ...func(*UTicker)) *UTicker {

	t := &UTicker{
		C:              make(chan time.Time),
		Duration:       1 * time.Second,
		ImmediateStart: false,
	}

	for _, option := range options {
		option(t)
	}

	t.ticker = time.NewTicker(t.Duration)

	go func() {
		if t.ImmediateStart {
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

	t.ticker.Stop()
	close(t.C)

}

func (t *UTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	t.ticker.Reset(d)
}

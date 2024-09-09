package uticker

import (
	"time"
)

type UTicker struct {
	C              chan time.Time
	Duration       time.Duration
	ImmediateStart bool
	NextTick       func() time.Duration
	ticker         *time.Ticker
	counter        uint64
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

func WithExponentialBackoff(e int) func(*UTicker) {
	return func(t *UTicker) {
		t.NextTick = func() time.Duration {
			return t.Duration * time.Duration(e)
		}
	}
}

func WithExponentialBackoffCapped(e int, max int) func(*UTicker) {
	return func(t *UTicker) {
		t.NextTick = func() time.Duration {
			if t.counter > uint64(max) {
				return t.Duration
			} else {
				return t.Duration * time.Duration(e)
			}
		}
	}
}

func WithRampCapped(e int, max int) func(*UTicker) {
	return func(t *UTicker) {
		t.NextTick = func() time.Duration {
			if t.counter > uint64(max) {
				return t.Duration
			} else {
				return t.Duration / time.Duration(e)
			}
		}
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
			tick(t)
		}
		for {
			select {
			case <-t.ticker.C:
				tick(t)
				if t.NextTick != nil {
					calculateNextTick(t)
				}
			}
		}
	}()
	return t
}

func calculateNextTick(t *UTicker) {
	t1 := t.NextTick()
	t.Reset(t1)
	t.Duration = t1
}

func tick(t *UTicker) {
	t.C <- time.Now()
	t.counter++
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

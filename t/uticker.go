package t

import (
	"github.com/gorhill/cronexpr"
	"math/rand/v2"
	"time"
)

type UTicker struct {
	C              chan time.Time
	frequency      time.Duration
	immediateStart bool
	nextTick       func() time.Duration
	ticker         *time.Ticker
	counter        uint64
}

func WithCronExpression(cron string) func(*UTicker) {

	return func(t *UTicker) {

		now := time.Now()
		nextTime := cronexpr.MustParse(cron).Next(now)
		duration := nextTime.Sub(now)
		t.frequency = duration
		t.immediateStart = false

		t.nextTick = func() time.Duration {
			now := time.Now()
			nextTime := cronexpr.MustParse(cron).Next(now)
			duration := nextTime.Sub(now)
			return duration
		}
	}
}

func WithImmediateStart(b bool) func(*UTicker) {
	return func(t *UTicker) {
		t.immediateStart = b
	}
}

func WithFrequency(d time.Duration) func(*UTicker) {
	if d <= 0 {
		panic("non-positive frequency for NewTicker")
	}
	return func(t *UTicker) {
		t.frequency = d
	}
}

func WithExponentialBackoff(e int) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			return t.frequency * time.Duration(e)
		}
	}
}

func WithExponentialBackoffCapped(e int, max int) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			if t.counter > uint64(max) {
				return t.frequency
			} else {
				return t.frequency * time.Duration(e)
			}
		}
	}
}

func WithRampCapped(e int, max int) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			if t.counter > uint64(max) {
				return t.frequency
			} else {
				return t.frequency / time.Duration(e)
			}
		}
	}
}

func WithDeviation(percentage float64) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			deviation := t.frequency * time.Duration(percentage)
			return t.frequency + deviation
		}
	}
}

func WithAnotherDurationWithGivenProbability(duration time.Duration, probability float64) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			if rand.Float64() < probability {
				return t.frequency
			} else {
				return duration
			}
		}
	}
}

func WithRandomTickIn(duration time.Duration) func(*UTicker) {
	return func(t *UTicker) {
		t.nextTick = func() time.Duration {
			d := rand.Float64() * float64(duration.Milliseconds())
			return time.Duration(d) * time.Millisecond
		}
	}
}

func NewUTicker(options ...func(*UTicker)) *UTicker {

	t := &UTicker{
		C:              make(chan time.Time),
		frequency:      1 * time.Second,
		immediateStart: false,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

func (t *UTicker) run() {

	if t.immediateStart {
		t.tick()
	}
	for {
		select {
		case <-t.ticker.C:
			t.tick()
			if t.nextTick != nil {
				t.calculateNextTick()
			}
		}
	}
}

func (t *UTicker) calculateNextTick() {
	t1 := t.nextTick()
	t.Reset(t1)
	t.frequency = t1
}

func (t *UTicker) tick() {
	t.C <- time.Now()
	t.counter++
}

func (t *UTicker) Stop() {
	t.ticker.Stop()
	// TODO: golang ticker doc says not to close the channel
	// Stop turns off a ticker. After Stop, no more ticks will be sent.
	// Stop does not close the channel, to prevent a concurrent goroutine
	// reading from the channel from seeing an erroneous "tick".
	//close(t.C)
}

func (t *UTicker) Start() {
	t.ticker = time.NewTicker(t.frequency)
	go t.run()
}

func (t *UTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for ticker.Reset")
	}
	t.ticker.Reset(d)
}

package uticker

import (
	"time"
)

type UTicker struct {
	C      chan time.Time
	ticker time.Ticker
}

func NewUTicker(d time.Duration) *UTicker {
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	t := &UTicker{
		C:      make(chan time.Time, 1),
		ticker: *time.NewTicker(d),
	}
	go func() {
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
	close(t.C)
	//t.ticker.Stop()
}

func (t *UTicker) Reset(d time.Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	t.ticker.Reset(d)
}

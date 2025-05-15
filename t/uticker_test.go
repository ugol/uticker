package t_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ugol/uticker/t"
	"testing"
	"time"
)

func TestSimpleTicker(test *testing.T) {
	ticker := t.NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Simple ticker at 1s", 5*time.Second, 5, test)
}

func TestImmediateStartTicker(test *testing.T) {
	ticker := t.NewUTicker(t.WithImmediateStart(true))
	defer ticker.Stop()
	runExample(ticker, "Immediate start ticker at 1s", 5*time.Second, 6, test)
}

func TestImmediateStartTickerWithFrequency(test *testing.T) {
	ticker := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithFrequency(100*time.Millisecond),
	)
	defer ticker.Stop()
	runExample(ticker, "Immediate start ticker at 100ms", 5*time.Second, 70, test)
}
func runExample(ticker *t.UTicker, msg string, d time.Duration, expected int, test *testing.T) {
	stop := make(chan struct{})
	run(ticker, msg, stop)

	timer := time.NewTimer(d)
	defer timer.Stop()

	assert.Equal(test, expected, int(ticker.Counter))

}

func run(ticker *t.UTicker, msg string, testStop chan struct{}) {
	fmt.Println(msg)
	ticker.Start()
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case tickTime := <-ticker.C:
				fmt.Println("Tick at", tickTime)
			case <-stop:
				return
			case <-testStop:
				close(stop)
				return
			}
		}
	}()
}

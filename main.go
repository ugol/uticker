package main

import (
	"fmt"
	. "github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Normal ticker at 1s", 3*time.Second)

	ticker1 := NewUTicker(WithImmediateStart())
	defer ticker1.Stop()
	runExample(ticker1, "Immediate start ticker at 1s", 5*time.Second)

	ticker2 := NewUTicker(
		WithImmediateStart(),
		WithDuration(100*time.Millisecond),
	)
	defer ticker2.Stop()
	runExample(ticker2, "Immediate start ticker at 100ms", 5*time.Second)

	ticker3 := NewUTicker(
		WithImmediateStart(),
		WithDuration(100*time.Millisecond),
		WithExponentialBackoff(2),
	)
	defer ticker3.Stop()
	runExample(ticker3, "Immediate start ticker at 100ms with Exponential backoff", 3*time.Second)

	ticker4 := NewUTicker(
		WithImmediateStart(),
		WithDuration(100*time.Millisecond),
		WithExponentialBackoffCapped(2, 3),
	)
	defer ticker4.Stop()
	runExample(ticker4, "Immediate start ticker at 100ms with Exponential backoff and cap", 3*time.Second)

	ticker5 := NewUTicker(
		WithImmediateStart(),
		WithDuration(5*time.Second),
		WithRampCapped(2, 10),
	)
	defer ticker5.Stop()
	runExample(ticker5, "Immediate start ticker at 5s with ramp (halfing each tick) and cap", 10*time.Second)

}

func runExample(ticker *UTicker, msg string, d time.Duration) {

	fmt.Println(msg)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()
	time.Sleep(d)
	done <- true
}

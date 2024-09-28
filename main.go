package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	ticker := NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Normal ticker at 1s", 3*time.Second)

	ticker1 := NewUTicker(
		WithImmediateStart(true),
	)
	defer ticker1.Stop()
	runExample(ticker1, "Immediate start ticker at 1s", 5*time.Second)

	ticker2 := NewUTicker(
		WithImmediateStart(true),
		WithFrequency(100*time.Millisecond),
	)
	defer ticker2.Stop()
	runExample(ticker2, "Immediate start ticker at 100ms", 5*time.Second)

	ticker3 := NewUTicker(
		WithImmediateStart(true),
		WithFrequency(100*time.Millisecond),
		WithExponentialBackoff(2),
	)
	defer ticker3.Stop()
	runExample(ticker3, "Immediate start ticker at 100ms with Exponential backoff", 3*time.Second)

	ticker4 := NewUTicker(
		WithImmediateStart(true),
		WithFrequency(100*time.Millisecond),
		WithExponentialBackoffCapped(2, 3),
	)
	defer ticker4.Stop()
	runExample(ticker4, "Immediate start ticker at 100ms with Exponential backoff and cap", 3*time.Second)

	ticker5 := NewUTicker(
		WithImmediateStart(true),
		WithFrequency(5*time.Second),
		WithRampCapped(2, 10),
	)
	defer ticker5.Stop()
	runExample(ticker5, "Immediate start ticker at 5s with ramp (halfing each tick) and cap", 10*time.Second)

	ticker6 := NewUTicker(
		WithImmediateStart(true),
		WithRandomTickIn(1*time.Second),
	)
	defer ticker6.Stop()
	runExample(ticker6, "Immediate start ticker with random tick in 1S", 5*time.Second)

}

func runExample(ticker *UTicker, msg string, d time.Duration) {
	fmt.Println(msg)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			case <-stop:
				return
			}
		}
	}()

	time.AfterFunc(d, func() {
		close(stop)
	})
	wg.Wait()
}

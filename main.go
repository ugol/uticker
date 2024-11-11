package main

import (
	"fmt"
	"github.com/ugol/uticker/t"
	"sync"
	"time"
)

func main() {

	ticker := t.NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Normal ticker at 1s", 3*time.Second)

	ticker1 := t.NewUTicker(
		t.WithImmediateStart(true),
	)
	defer ticker1.Stop()
	runExample(ticker1, "Immediate start ticker at 1s", 5*time.Second)

	ticker2 := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithFrequency(100*time.Millisecond),
	)
	defer ticker2.Stop()
	runExample(ticker2, "Immediate start ticker at 100ms", 5*time.Second)

	ticker3 := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithFrequency(100*time.Millisecond),
		t.WithExponentialBackoff(2),
	)
	defer ticker3.Stop()
	runExample(ticker3, "Immediate start ticker at 100ms with Exponential backoff", 3*time.Second)

	ticker4 := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithFrequency(100*time.Millisecond),
		t.WithExponentialBackoffCapped(2, 3),
	)
	defer ticker4.Stop()
	runExample(ticker4, "Immediate start ticker at 100ms with Exponential backoff and cap", 3*time.Second)

	ticker5 := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithFrequency(5*time.Second),
		t.WithRampCapped(2, 10),
	)
	defer ticker5.Stop()
	runExample(ticker5, "Immediate start ticker at 5s with ramp (halfing each tick) and cap", 10*time.Second)

	ticker6 := t.NewUTicker(
		t.WithImmediateStart(true),
		t.WithRandomTickIn(1*time.Second),
	)
	defer ticker6.Stop()
	runExample(ticker6, "Immediate start ticker with random tick in 1S", 5*time.Second)

	cron := "*/3 * * * * * *"
	ticker7 := t.NewUTicker(
		t.WithCronExpression(cron),
	)
	defer ticker7.Stop()
	runExample(ticker7, "Ticker with cron expression: "+cron, 30*time.Second)

}

func runExample(ticker *t.UTicker, msg string, d time.Duration) {
	fmt.Println(msg)
	ticker.Start()
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

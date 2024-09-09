package main

import (
	"fmt"
	. "github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Normal ticker at 1s")

	ticker1 := NewUTicker(WithImmediateStart())
	defer ticker1.Stop()
	runExample(ticker1, "Immediate start ticker at 1s")

	ticker2 := NewUTicker(
		WithImmediateStart(),
		WithDuration(100*time.Millisecond),
	)
	defer ticker2.Stop()
	runExample(ticker2, "Immediate start ticker at 100ms")

	ticker3 := NewUTicker(
		WithImmediateStart(),
		WithDuration(100*time.Millisecond),
		Exponential(),
	)
	defer ticker3.Stop()
	runExample(ticker3, "Immediate start ticker at 100ms with Exponential backoff")

}

func runExample(ticker *UTicker, msg string) {

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
	time.Sleep(3 * time.Second)
	done <- true
}

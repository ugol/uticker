package main

import (
	"fmt"
	. "github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := NewUTicker()
	runExample(ticker, "Normal ticker")

	ticker1 := NewUTicker(WithImmediateStart)
	runExample(ticker1, "Immediate start ticker")

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
	time.Sleep(6 * time.Second)
	ticker.Stop()
	done <- true
}

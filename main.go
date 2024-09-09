package main

import (
	"fmt"
	. "github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := NewUTicker()
	defer ticker.Stop()
	runExample(ticker, "Normal ticker")

	ticker1 := NewUTicker(WithImmediateStart)
	defer ticker1.Stop()
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
	done <- true
}

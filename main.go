package main

import (
	"fmt"
	"github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := uticker.NewUTicker()
	runExample(ticker, "Normal ticker")

	ticker1 := uticker.NewUTicker(uticker.WithImmediateStart)
	runExample(ticker1, "Immediate start ticker")

}

func runExample(ticker *uticker.UTicker, msg string) {

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
	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true
}

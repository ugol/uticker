package main

import (
	"fmt"
	"github.com/ugol/uticker/uticker"
	"time"
)

func main() {

	ticker := uticker.NewUTicker(1 * time.Second)
	go func() {
		for ; ; _ = <-ticker.C {
			fmt.Println("Simple ticker")
		}
	}()

	time.Sleep(5 * time.Second)
	ticker.Stop()

}

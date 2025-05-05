package t_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ugol/uticker/t"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestSimpleTicker(test *testing.T) {
	synctest.Run(func() {
		ticker := t.NewUTicker()
		// Create stop channel before starting
		stop := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			runExample(ticker, "Normal ticker at 1s", 30500*time.Millisecond)
		}()

		// Ensure cleanup happens
		defer func() {
			close(stop)
			ticker.Stop()
			wg.Wait()
		}()

		// Wait for completion or timeout
		timer := time.NewTimer(31 * time.Second)
		defer timer.Stop()

		select {
		case <-timer.C:
			return
		case <-stop:
			assert.Equal(test, 30, int(ticker.Counter))
		}
	})
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
			case tickTime := <-ticker.C:
				fmt.Println("Tick at", tickTime)
			case <-stop:
				return
			}
		}
	}()

	time.AfterFunc(d, func() {
		close(stop)
	})
	wg.Wait()

	fmt.Printf("Ticks: %d\n", ticker.Counter)
}

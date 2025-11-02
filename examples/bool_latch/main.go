package main

import (
	syngo_sync "github.com/synesissoftware/syngo/sync"

	"fmt"
	"os"
	"sync"
	"time"
)

func main() {

	latch := syngo_sync.NewBoolLatch()

	var wg sync.WaitGroup

	wg.Go(func() {

		for i := 0; ; i++ {

			if latch.Load() {

				fmt.Fprintf(os.Stdout, "detected latched after %d iteration(s)\n", i)
				break
			}
		}
	})

	wg.Go(func() {

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		<-ticker.C

		latch.Set()
	})

	wg.Wait()
}

package Queue

import (
	"testing"
	//"time"
	"fmt"
	"sync"
)

// / go test -v
func TestNqueue(t *testing.T) {
	qqqq := NewNQueue[int]()

	go qqqq.DequeueFunc(func(fff int, isClose bool) bool {
		fmt.Println(fff)

		return true
	})

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i < 100; i++ {
				qqqq.Enqueue(i)
			}
		}()

	}

	wg.Wait()
	qqqq.Close()
}

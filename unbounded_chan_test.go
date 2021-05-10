package chanx

import (
	"fmt"
	"sync"
	"testing"
)

func TestMakeUnboundedChan(t *testing.T) {
	in, out := MakeUnboundedChan(100)

	for i := 1; i < 200; i++ {
		in <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range out {
			count += v.(int64)
		}

		fmt.Println("read completed")
	}()

	for i := 200; i <= 1000; i++ {
		in <- int64(i)
	}
	close(in)

	wg.Wait()

	if count != 500500 {
		t.Fatalf("expected 500500 but got %d", count)
	}
}

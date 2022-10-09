package chanx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeUnboundedChan(t *testing.T) {
	ch := NewUnboundedChan[int64](100)

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range ch.Out {
			count += v
		}
	}()

	for i := 200; i <= 1000; i++ {
		ch.In <- int64(i)
	}
	close(ch.In)

	wg.Wait()

	if count != 500500 {
		t.Fatalf("expected 500500 but got %d", count)
	}
}

func TestMakeUnboundedChanSize(t *testing.T) {
	ch := NewUnboundedChanSize[int64](10, 50, 100)

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range ch.Out {
			count += v
		}
	}()

	for i := 200; i <= 1000; i++ {
		ch.In <- int64(i)
	}
	close(ch.In)

	wg.Wait()

	if count != 500500 {
		t.Fatalf("expected 500500 but got %d", count)
	}
}

func TestLen_DataRace(t *testing.T) {
	ch := NewUnboundedChan[int64](1)
	stop := make(chan bool)
	for i := 0; i < 100; i++ { // may tweak the number of iterations
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					ch.In <- 42
					<-ch.Out
				}
			}
		}()
	}

	for i := 0; i < 10000; i++ { // may tweak the number of iterations
		ch.Len()
	}
	close(stop)
}

func TestLen(t *testing.T) {
	ch := NewUnboundedChanSize[int64](10, 50, 100)

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	// wait ch processing in normal case
	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 50, len(ch.Out))
	assert.Equal(t, 199, ch.Len())
	assert.Equal(t, 149, ch.BufLen())

	for i := 0; i < 50; i++ {
		<-ch.Out
	}

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 50, len(ch.Out))
	assert.Equal(t, 149, ch.Len())
	assert.Equal(t, 99, ch.BufLen())

	for i := 0; i < 149; i++ {
		<-ch.Out
	}

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 0, len(ch.Out))
	assert.Equal(t, 0, ch.Len())
	assert.Equal(t, 0, ch.BufLen())
}

func TestCancel(t *testing.T) {
	ch := NewUnboundedChan[int64](1)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch.Done:
					return
				case ch.In <- 42:
				}
			}
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch.Done:
					return
				case _ = <-ch.Out:
				}
			}
		}()
	}

	time.Sleep(time.Second)
	ch.Cancel()
	wg.Wait()
}

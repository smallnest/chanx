package chanx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMakeUnboundedChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewUnboundedChan[int64](ctx, 100)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewUnboundedChanSize[int64](ctx, 10, 50, 100)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewUnboundedChan[int64](ctx, 1)
	stop := make(chan bool)
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ { // may tweak the number of iterations
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-stop:
					return
				default:
					ch.In <- 42
					<-ch.Out
				}
			}
		}(ctx)
	}

	for i := 0; i < 10000; i++ { // may tweak the number of iterations
		ch.Len()
	}
	close(stop)
	wg.Wait()
}

func TestLen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewUnboundedChanSize[int64](ctx, 10, 50, 100)

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

func TestGetDataWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewUnboundedChanSize[int64](ctx, 10, 50, 100)
	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}
}

package chanx

import (
	"sync/atomic"
)

// T defines interface{}, and will be used for generic type after go 1.18 is released.
type T interface{}

// UnboundedChan is an unbounded chan.
// In is used to write without blocking, which supports multiple writers.
// and Out is used to read, which supports multiple readers.
// You can close the in channel if you want.
type UnboundedChan struct {
	bufCount int64
	In       chan<- T    // channel for write
	Out      <-chan T    // channel for read
	buffer   *RingBuffer // buffer
}

// Len returns len of In plus len of Out plus len of buffer.
// It is not accurate and only for your evaluating approximate number of elements in this chan,
// see https://github.com/smallnest/chanx/issues/7.
func (c UnboundedChan) Len() int {
	return len(c.In) + c.BufLen() + len(c.Out)
}

// BufLen returns len of the buffer.
// It is not accurate and only for your evaluating approximate number of elements in this chan,
// see https://github.com/smallnest/chanx/issues/7.
func (c UnboundedChan) BufLen() int {
	return int(atomic.LoadInt64(&c.bufCount))
}

// NewUnboundedChan creates the unbounded chan.
// in is used to write without blocking, which supports multiple writers.
// and out is used to read, which supports multiple readers.
// You can close the in channel if you want.
func NewUnboundedChan(initCapacity int) *UnboundedChan {
	return NewUnboundedChanSize(initCapacity, initCapacity, initCapacity)
}

// NewUnboundedChanSize is like NewUnboundedChan but you can set initial capacity for In, Out, Buffer.
func NewUnboundedChanSize(initInCapacity, initOutCapacity, initBufCapacity int) *UnboundedChan {
	in := make(chan T, initInCapacity)
	out := make(chan T, initOutCapacity)
	ch := UnboundedChan{In: in, Out: out, buffer: NewRingBuffer(initBufCapacity)}

	go process(in, out, &ch)

	return &ch
}

func process(in, out chan T, ch *UnboundedChan) {
	for {
		if ch.buffer.IsEmpty() {
			val, unclosed := <-in
			if !unclosed {
				goto exit
			}
			ch.buffer.Write(val)
			atomic.AddInt64(&ch.bufCount, 1)
		}

		select {
		case val, unclosed := <-in:
			if !unclosed {
				goto exit
			}
			ch.buffer.Write(val)
			atomic.AddInt64(&ch.bufCount, 1)
		case out <- ch.buffer.Peek():
			ch.buffer.Pop()
			atomic.AddInt64(&ch.bufCount, -1)
			if ch.buffer.IsEmpty() && ch.buffer.size > ch.buffer.initialSize {
				ch.buffer.Reset()
				atomic.StoreInt64(&ch.bufCount, 0)
			}
		}
	}

exit:
	for !ch.buffer.IsEmpty() {
		out <- ch.buffer.Pop()
		atomic.AddInt64(&ch.bufCount, -1)
	}

	ch.buffer.Reset()
	atomic.StoreInt64(&ch.bufCount, 0)

	close(out)
}

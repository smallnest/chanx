package chanx

// // UnboundedChan is an unbounded chan.
// // In is used to write without blocking, which supports multiple writers.
// // and Out is used to read, wich supports multiple readers.
// // You can close the in channel if you want.
// type UnboundedChan[T any] struct {
// 	In     chan<- T // channel for write
// 	Out    <-chan T // channel for read
// 	buffer []T      // buffer
// }

// // Len returns len of Out plus len of buffer.
// func (c UnboundedChan) Len() int {
// 	return len(c.buffer) + len(c.Out)
// }

// // BufLen returns len of the buffer.
// func (c UnboundedChan) BufLen() int {
// 	return len(c.buffer)
// }

// // NewUnboundedChan creates the unbounded chan.
// // in is used to write without blocking, which supports multiple writers.
// // and out is used to read, wich supports multiple readers.
// // You can close the in channel if you want.
// func NewUnboundedChan(initCapacity int) UnboundedChan {
// 	in := make(chan T, initCapacity)
// 	out := make(chan T, initCapacity)
// 	ch := UnboundedChan{In: in, Out: out, buffer: make([]T, 0, initCapacity)}

// 	go func() {
// 		defer close(out)
// 	loop:
// 		for {
// 			val, ok := <-in
// 			if !ok { // in is closed
// 				break loop
// 			}

// 			// out is not full
// 			select {
// 			case out <- val:
// 				continue
// 			default:
// 			}

// 			// out is full
// 			ch.buffer = append(ch.buffer, val)
// 			for len(ch.buffer) > 0 {
// 				select {
// 				case val, ok := <-in:
// 					if !ok { // in is closed
// 						break loop
// 					}
// 					ch.buffer = append(ch.buffer, val)

// 				case out <- ch.buffer[0]:
// 					ch.buffer = ch.buffer[1:]
// 					if len(ch.buffer) == 0 { // after burst
// 						ch.buffer = make([]T, 0, initCapacity)
// 						ch.buffer = ch.buffer
// 					}
// 				}
// 			}
// 		}

// 		// drain
// 		for len(ch.buffer) > 0 {
// 			out <- ch.buffer[0]
// 			ch.buffer = ch.buffer[1:]
// 		}
// 	}()

// 	return ch
// }

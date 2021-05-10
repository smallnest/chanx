package chanx

// MakeUnboundedChan creates the unbounded chan.
// in is used to write without blocking, which supports multiple writers.
// and out is used to read, wich supports multiple readers.
// You can close the in channel if you want.
func MakeUnboundedChan(initCapacity int) (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{}, initCapacity)
	out := make(chan interface{}, initCapacity)

	go func() {
		defer close(out)
		buffer := make([]interface{}, 0, initCapacity)
	loop:
		for {
			val, ok := <-in
			if !ok { // in is closed
				break loop
			}

			// out is not full
			select {
			case out <- val:
				continue
			default:
			}

			// out is full
			buffer = append(buffer, val)
			for len(buffer) > 0 {
				select {
				case val, ok := <-in:
					if !ok { // in is closed
						break loop
					}
					buffer = append(buffer, val)

				case out <- buffer[0]:
					buffer = buffer[1:]
					if len(buffer) == 0 { // after burst
						buffer = make([]interface{}, 0, initCapacity)
					}
				}
			}
		}

		// drain
		for len(buffer) > 0 {
			out <- buffer[0]
			buffer = buffer[1:]
		}
	}()

	return in, out
}

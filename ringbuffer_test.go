package chanx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	rb := NewRingBuffer(10)
	v, err := rb.Read()
	assert.Nil(t, v)
	assert.Error(t, err, ErrIsEmpty)

	write := 0
	read := 0

	// write one and read it
	rb.Write(0)
	v, err = rb.Read()
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
	assert.Equal(t, 1, rb.r)
	assert.Equal(t, 1, rb.w)
	assert.True(t, rb.IsEmpty())

	// then write 10
	for i := 0; i < 9; i++ {
		rb.Write(i)
		write += i
	}
	assert.Equal(t, 10, rb.Capacity())
	assert.Equal(t, 9, rb.Len())

	// write one more, the buffer is full so it grows
	rb.Write(10)
	write += 10
	assert.Equal(t, 20, rb.Capacity())
	assert.Equal(t, 10, rb.Len())

	for i := 0; i < 90; i++ {
		rb.Write(i)
		write += i
	}

	assert.Equal(t, 160, rb.Capacity())
	assert.Equal(t, 100, rb.Len())

	for {
		v, err := rb.Read()
		if err == ErrIsEmpty {
			break
		}

		read += v.(int)
	}

	assert.Equal(t, write, read)

	rb.Reset()
	assert.Equal(t, 10, rb.Capacity())
	assert.Equal(t, 0, rb.Len())
	assert.True(t, rb.IsEmpty())
}

func TestRingBuffer_One(t *testing.T) {
	rb := NewRingBuffer(1)
	v, err := rb.Read()
	assert.Nil(t, v)
	assert.Error(t, err, ErrIsEmpty)

	write := 0
	read := 0

	// write one and read it
	rb.Write(0)
	v, err = rb.Read()
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
	assert.Equal(t, 1, rb.r)
	assert.Equal(t, 1, rb.w)
	assert.True(t, rb.IsEmpty())

	// then write 10
	for i := 0; i < 9; i++ {
		rb.Write(i)
		write += i
	}
	assert.Equal(t, 16, rb.Capacity())
	assert.Equal(t, 9, rb.Len())

	// write one more, the buffer is full so it grows
	rb.Write(10)
	write += 10
	assert.Equal(t, 16, rb.Capacity())
	assert.Equal(t, 10, rb.Len())

	for i := 0; i < 90; i++ {
		rb.Write(i)
		write += i
	}

	assert.Equal(t, 128, rb.Capacity())
	assert.Equal(t, 100, rb.Len())

	for {
		v, err := rb.Read()
		if err == ErrIsEmpty {
			break
		}

		read += v.(int)
	}

	assert.Equal(t, write, read)

	rb.Reset()
	assert.Equal(t, 2, rb.Capacity())
	assert.Equal(t, 0, rb.Len())
	assert.True(t, rb.IsEmpty())
}

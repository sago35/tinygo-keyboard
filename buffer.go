//go:build tinygo

package keyboard

import (
	"runtime/volatile"
)

// RingBuffer is ring buffer implementation inspired by post at
// https://www.embeddedrelated.com/showthread/comp.arch.embedded/77084-1.php
type RingBuffer[T any] struct {
	buffer []T
	head   volatile.Register8
	tail   volatile.Register8
}

// NewRingBuffer returns a new ring buffer.
func NewRingBuffer[T any](buf []T) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: buf,
	}
}

// Used returns how many bytes in buffer have been used.
func (rb *RingBuffer[T]) Used() uint8 {
	return uint8(rb.head.Get() - rb.tail.Get())
}

// Put stores a byte in the buffer. If the buffer is already
// full, the method will return false.
func (rb *RingBuffer[T]) Put(val T) bool {
	if rb.Used() != uint8(len(rb.buffer)) {
		rb.head.Set(rb.head.Get() + 1)
		rb.buffer[rb.head.Get()%uint8(len(rb.buffer))] = val
		return true
	}
	return false
}

// Get returns a byte from the buffer. If the buffer is empty,
// the method will return a false as the second value.
func (rb *RingBuffer[T]) Get() (T, bool) {
	if rb.Used() != 0 {
		rb.tail.Set(rb.tail.Get() + 1)
		return rb.buffer[rb.tail.Get()%uint8(len(rb.buffer))], true
	}
	var ret T
	return ret, false
}

// Peek peeks a byte from the buffer. If the buffer is empty,
// the method will return a false as the second value.
func (rb *RingBuffer[T]) Peek() (T, bool) {
	if rb.Used() != 0 {
		return rb.buffer[(rb.tail.Get()+1)%uint8(len(rb.buffer))], true
	}
	var ret T
	return ret, false
}

// Clear resets the head and tail pointer to zero.
func (rb *RingBuffer[T]) Clear() {
	rb.head.Set(0)
	rb.tail.Set(0)
}

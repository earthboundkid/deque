//go:build go1.23

package deque

import (
	"iter"
	"slices"
)

// Slice returns a slice with a copy of the deque.
func (d *Deque[T]) Slice() []T {
	s := make([]T, d.Len())
	for i, v := range d.All() {
		s[i] = v
	}
	return s
}

// All returns a sequence yielding each index and value in the deque.
func (d *Deque[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range d.Len() {
			v, ok := d.At(i)
			if !ok || !yield(i, v) {
				return
			}
		}
	}
}

// Reverse returns a sequence yielding each index and value in the deque in reverse order.
func (d *Deque[T]) Reverse() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := d.Len() - 1; i >= 0; i-- {
			v, ok := d.At(i)
			if !ok || !yield(i, v) {
				return
			}
		}
	}
}

// PushBackSeq adds all items in seq to the back of the deque.
func (d *Deque[T]) PushBackSeq(seq iter.Seq[T]) {
	for v := range seq {
		d.PushBack(v)
	}
}

// PushBackSlice adds all items in s to the back of the deque.
func (d *Deque[T]) PushBackSlice(s []T) {
	d.Grow(len(s))
	d.PushBackSeq(slices.Values(s))
}

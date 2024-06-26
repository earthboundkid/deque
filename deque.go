// Package deque provides a type-safe, slice-backed, double-ended queue.
//
// See https://en.wikipedia.org/wiki/Double-ended_queue
package deque

import (
	"cmp"
	"fmt"
	"strings"
)

// Deque is a double-ended queue. It is not concurrency safe.
type Deque[T any] struct {
	len, head int
	backing   []T
}

// Make creates a deque with a prereserved capacity.
func Make[T any](cap int) *Deque[T] {
	var d Deque[T]
	d.Grow(cap)
	return &d
}

// Of constructs a new Deque.
func Of[T any](items ...T) *Deque[T] {
	var d Deque[T]
	d.Grow(len(items))
	for i := range items {
		d.PushBack(items[i])
	}
	return &d
}

// Grow increases the deque's capacity, if necessary, to guarantee space for another n elements.
// After Grow(n), at least n elements can be appended to the deque without another allocation.
// If n is negative, Grow panics.
func (d *Deque[T]) Grow(n int) {
	if n < 0 {
		panic("argument to Grow must be positive")
	}
	if d.Cap()-d.len >= n {
		return
	}
	// using append to get amortized growth
	grown := append(d.backing, make([]T, n)...)
	grown = grown[:cap(grown)]
	d.copy(grown)
}

// Len returns the current length of the deque.
func (d *Deque[T]) Len() int {
	return d.len
}

// Cap returns the total current capacity of the deque.
func (d *Deque[T]) Cap() int {
	return len(d.backing)
}

// PushFront adds a new value v to the front of the deque.
func (d *Deque[T]) PushFront(v T) {
	d.Grow(1)
	d.len++
	d.head--
	if d.head < 0 {
		d.head = d.Cap() - 1
	}
	d.backing[d.head] = v
}

func (d *Deque[T]) copy(dst []T) {
	front, back := d.frontback()
	n := copy(dst, front)
	copy(dst[n:], back)
	d.head = 0
	d.backing = dst
}

// Clip removes unused capacity from the deque.
func (d *Deque[T]) Clip() {
	if d.Cap() == d.Len() {
		return
	}
	d.copy(make([]T, d.Len()))
}

// Front returns the first value of the deque,
// if any.
func (d *Deque[T]) Front() (v T, ok bool) {
	if p := d.at(0); p != nil {
		return *p, true
	}
	return
}

func (d *Deque[T]) tail() *T {
	return d.at(d.len - 1)
}

func (d *Deque[T]) at(n int) *T {
	if n < 0 || n > d.len-1 {
		return nil
	}
	return &d.backing[(d.head+n)%d.Cap()]
}

// At returns the zero indexed nth item of the deque, if any.
func (d *Deque[T]) At(n int) (t T, ok bool) {
	if p := d.at(n); p != nil {
		return *p, true
	}
	return
}

// Back returns the last value of the deque,
// if any.
func (d *Deque[T]) Back() (t T, ok bool) {
	if p := d.tail(); p != nil {
		return *p, true
	}
	return
}

// PushBack adds new value v to the end of the deque.
func (d *Deque[T]) PushBack(v T) {
	d.Grow(1)
	d.len++
	*d.tail() = v
}

// PushBackSlice adds all items in s to the back of the deque.
func (d *Deque[T]) PushBackSlice(s []T) {
	d.Grow(len(s))
	for _, t := range s {
		d.len++
		*d.tail() = t
	}
}

// RemoveFront removes and returns the front of the deque,
// if any.
func (d *Deque[T]) RemoveFront() (t T, ok bool) {
	if d.len < 1 {
		return
	}
	head, _ := d.Front()
	d.head++
	if d.head >= d.Cap() {
		d.head = 0
	}
	d.len--
	return head, true
}

// RemoveBack removes and returns the back of the deque,
// if any.
func (d *Deque[T]) RemoveBack() (t T, ok bool) {
	tail, ok := d.Back()
	if !ok {
		return
	}
	d.len--
	return tail, true
}

func (d *Deque[T]) frontback() (front, back []T) {
	end := min(d.head+d.len, len(d.backing))
	front = d.backing[d.head:end]
	rest := d.len - (end - d.head)
	back = d.backing[:rest]
	return
}

// String implements fmt.Stringer.
func (d *Deque[T]) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Deque{ len: %d, cap: %d, items: [", d.Len(), d.Cap())
	i := 0
	front, back := d.frontback()
	for _, slice := range [][]T{front, back} {
		for _, item := range slice {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprint(&buf, item)
			i++
		}
	}
	buf.WriteString("]}")
	return buf.String()
}

// Swap swaps the elements with indexes i and j.
func (d *Deque[T]) Swap(i, j int) {
	if i > d.len {
		panic("i out of bounds")
	}
	if j > d.len {
		panic("j out of bounds")
	}
	*d.at(i), *d.at(j) = *d.at(j), *d.at(i)
}

// Sortable is a deque that can be sorted with sort.Sort.
type Sortable[T cmp.Ordered] struct {
	*Deque[T]
}

// Less implements sort.Interface.
func (sd Sortable[T]) Less(i, j int) bool {
	if i > sd.len {
		panic("i out of bounds")
	}
	if j > sd.len {
		panic("j out of bounds")
	}
	return *sd.at(i) < *sd.at(j)
}

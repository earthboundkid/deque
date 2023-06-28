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
		d.PushTail(items[i])
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

// Cap returns the unused capacity of the deque.
func (d *Deque[T]) Cap() int {
	return cap(d.backing)
}

// PushHead adds t to the head of the deque.
func (d *Deque[T]) PushHead(t T) {
	d.Grow(1)
	d.len++
	d.head--
	if d.head < 0 {
		d.head = d.Cap() - 1
	}
	d.backing[d.head] = t
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

// Head returns the head of the deque, if any.
func (d *Deque[T]) Head() (t T, ok bool) {
	if d.len < 1 {
		return
	}
	return d.backing[d.head], true
}

func (d *Deque[T]) tail() int {
	return d.at(d.len - 1)
}

func (d *Deque[T]) at(n int) int {
	if d.len < 1 {
		return -1
	}
	return (d.head + n) % d.Cap()
}

// At returns the zero indexed nth item of the deque, if any.
func (d *Deque[T]) At(n int) (t T, ok bool) {
	if n < 0 || n > d.len-1 {
		return
	}
	return d.backing[d.at(n)], true
}

// Tail returns the tail of the deque, if any.
func (d *Deque[T]) Tail() (t T, ok bool) {
	if d.len < 1 {
		return
	}
	return d.backing[d.tail()], true
}

// PushTail adds t to the tail of the deque.
func (d *Deque[T]) PushTail(t T) {
	d.Grow(1)
	d.len++
	d.backing[d.tail()] = t
}

// Append pushes all items to the tail of the deque.
func (d *Deque[T]) Append(ts ...T) {
	d.Grow(len(ts))
	for _, t := range ts {
		d.len++
		d.backing[d.tail()] = t
	}
}

// PopHead returns and removes the head of the deque, if any.
func (d *Deque[T]) PopHead() (t T, ok bool) {
	if d.len < 1 {
		return
	}
	head, _ := d.Head()
	d.head++
	if d.head >= d.Cap() {
		d.head = 0
	}
	d.len--
	return head, true
}

// PopTail returns and removes the tail of the deque, if any.
func (d *Deque[T]) PopTail() (t T, ok bool) {
	tail, ok := d.Tail()
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

// Slice returns a slice with a copy of the deque.
func (d *Deque[T]) Slice() []T {
	front, back := d.frontback()
	return append(append(([]T)(nil), front...), back...)
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
func (d Deque[T]) Swap(i, j int) {
	d.backing[d.at(i)], d.backing[d.at(j)] = d.backing[d.at(j)], d.backing[d.at(i)]
}

// Sortable is a deque that can be sorted with sort.Sort.
type Sortable[T cmp.Ordered] struct {
	*Deque[T]
}

// Less implements sort.Interface.
func (sd Sortable[T]) Less(i, j int) bool {
	return sd.backing[sd.at(i)] < sd.backing[sd.at(j)]
}

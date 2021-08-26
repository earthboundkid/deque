package deque

import (
	"fmt"
	"strings"
)

type Deque[T any] struct {
	len, head int
	backing   []T
}

func Make[T any](cap int) *Deque[T] {
	var d Deque[T]
	d.Grow(cap)
	return &d
}

func Of[T any](items ...T) *Deque[T] {
	var d Deque[T]
	d.Grow(len(items))
	for i := range items {
		d.PushTail(items[i])
	}
	return &d
}

func (d *Deque[T]) Grow(n int) {
	if n < 0 {
		panic("argument to Grow must be positive. did you want Shrink?")
	}
	if d.Cap()-d.len >= n {
		return
	}
	// using append to get amortized growth
	grown := append(d.backing, make([]T, n)...)
	grown = grown[:cap(grown)]
	d.copy(grown)
}

func (d *Deque[T]) Len() int {
	return d.len
}

func (d *Deque[T]) Cap() int {
	return cap(d.backing)
}

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

func (d *Deque[T]) Shrink() {
	if d.Cap() == d.Len() {
		return
	}
	d.copy(make([]T, d.Len()))
}

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

func (d *Deque[T]) At(n int) (t T, ok bool) {
	if n < 0 || n > d.len-1 {
		return
	}
	return d.backing[d.at(n)], true
}

func (d *Deque[T]) Tail() (t T, ok bool) {
	if d.len < 1 {
		return
	}
	return d.backing[d.tail()], true
}

func (d *Deque[T]) PushTail(t T) {
	d.Grow(1)
	d.len++
	d.backing[d.tail()] = t
}

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

func (d *Deque[T]) PopTail() (t T, ok bool) {
	tail, ok := d.Tail()
	if !ok {
		return
	}
	d.len--
	return tail, true
}

func (d *Deque[T]) frontback() (front, back []T) {
	end := d.head + d.len
	if end > len(d.backing) {
		end = len(d.backing)
	}
	front = d.backing[d.head:end]
	rest := d.len - (end - d.head)
	back = d.backing[:rest]
	return
}

func (d *Deque[T]) Items() []T {
	front, back := d.frontback()
	return append(append(([]T)(nil), front...), back...)
}

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

type ordered interface {
	int | string // ...
}

type Sortable[T ordered] struct {
	*Deque[T]
}

func (sd Sortable[T]) Less(i, j int) bool {
	return sd.backing[sd.at(i)] < sd.backing[sd.at(j)]
}

func (sd Sortable[T]) Swap(i, j int) {
	sd.backing[sd.at(i)], sd.backing[sd.at(j)] = sd.backing[sd.at(j)], sd.backing[sd.at(i)]
}

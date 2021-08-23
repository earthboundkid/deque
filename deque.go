package deque

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
	d.head -= 1
	if d.head < 0 {
		d.head = d.Cap() - 1
	}
	d.backing[d.head] = t
}

func (d *Deque[T]) copy(dst []T) {
	n := copy(dst, d.backing[d.head:])
	copy(dst[n:], d.backing[:d.head])
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
	if d.len < 1 {
		return -1
	}
	return (d.head + d.len - 1) % d.Cap()
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

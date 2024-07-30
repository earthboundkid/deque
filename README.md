# deque [![GoDoc](https://godoc.org/github.com/earthboundkid/deque?status.svg)](https://pkg.go.dev/github.com/earthboundkid/deque/v2) [![Coverage Status](https://coveralls.io/repos/github/earthboundkid/deque/badge.svg)](https://coveralls.io/github/earthboundkid/deque)
Deque is deque container using Go generics.

[Wikipedia says](https://en.wikipedia.org/wiki/Double-ended_queue):

> In computer science, a **double-ended queue** (abbreviated to **deque**, pronounced _deck_, like "cheque") is an abstract data type that generalizes a queue, for which elements can be added to or removed from either the front (head) or back (tail).

## Usage

```
// Make a new deque of ints
d := deque.Of(9, 8, 7, 6)

// Sort it
sort.Sort(deque.Sortable[int]{d})
// d is 6, 7, 8, 9

// Add 5, 4, 3, 2, 1 to the front
for i := 5; i > 0; i-- {
    d.PushFront(i)
}

// Deque{ len: 9, cap: 16, items: [1, 2, 3, 4, 5, 6, 7, 8, 9]}
fmt.Println(d)

// Now reverse loop through items
// Prints 9 8 7 6 5 4 3 2 1
for _, n := range d.Reverse() {
    fmt.Print(n, " ")
}
fmt.Println()
```

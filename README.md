# deque
Generic deque container


```
// Make a new deque of ints
d := deque.Of(9, 8, 7, 6)

// Sort it
sort.Sort(deque.Sortable[int]{d})
// d is 6, 7, 8, 9

// Add 5, 4, 3, 2, 1 to the front
for i := 5; i > 0; i-- {
    d.PushHead(i)
}

// Deque{ len: 9, cap: 16, items: [1, 2, 3, 4, 5, 6, 7, 8, 9]}
fmt.Println(d)

// Now pop items off the tail
// Prints 9 8 7 6 5 4 3 2 1
for {
    n, ok := d.PopTail()
    if !ok {
        break
    }
    fmt.Print(n, " ")
}
fmt.Println()
```

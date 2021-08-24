package deque

import (
	"fmt"
)

func ExampleDeque() {
	d := Of(6, 7, 8, 9)
	for i := 5; i > 0; i-- {
		d.PushHead(i)
	}
	fmt.Println(d)
	for {
		n, ok := d.PopTail()
		if !ok {
			break
		}
		fmt.Print(n, " ")
	}
	fmt.Println()
	// Output:
	// Deque{ len: 9, cap: 16, items: [1, 2, 3, 4, 5, 6, 7, 8, 9]}
	// 9 8 7 6 5 4 3 2 1
}

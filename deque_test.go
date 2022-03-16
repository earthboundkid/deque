package deque_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/carlmjohnson/deque"
)

func ExampleDeque() {
	// Make a new deque
	d := deque.Of(9, 8, 7, 6)
	// Sort it
	sort.Sort(deque.Sortable[int]{d})
	// Add 5, 4, 3, 2, 1 to the front
	for i := 5; i > 0; i-- {
		d.PushHead(i)
	}
	fmt.Println(d)
	// Now pop items off the tail
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

func dequeLang(t *testing.T, in string) {
	q := deque.Make[int](0)
	qlen := 0
	mincap := 0
	i := 0
	for _, c := range in {
		switch c {
		case '+':
			q.PushHead(i)
			qlen++
			if mincap < qlen {
				mincap++
			}
		case '*':
			q.PushTail(i)
			qlen++
			if mincap < qlen {
				mincap++
			}
		case '-':
			q.PopHead()
			qlen--
			if qlen < 0 {
				qlen = 0
			}
		case '/':
			q.PopTail()
			qlen--
			if qlen < 0 {
				qlen = 0
			}
		case '0':
			q.Clip()
			mincap = qlen
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n := int(c) - '0'
			q.Grow(n)
			if newcap := qlen + n; newcap > mincap {
				mincap = newcap
			}
		case 'A', 'B', 'C', 'D', 'E':
			n := int(c) - 'A'
			var s []int
			for j := 0; j < n; j++ {
				s = append(s, i)
				i++
			}
			q.Append(s...)
			qlen += n
			if mincap < qlen {
				mincap = qlen
			}

		}
		i++
	}
	if q.Len() != qlen {
		t.Errorf("%s bad len %d != %d", in, q.Len(), qlen)
	}
	if q.Cap() < mincap {
		t.Errorf("%s: bad cap %d < %d", in, q.Cap(), mincap)
	}
	seen := make(map[int]bool)
	s := q.Slice()
	if len(s) != q.Len() {
		t.Fatalf("slice has bad contents: %v != %v", s, q)
	}
	for i := 0; i < q.Len(); i++ {
		n, _ := q.At(i)
		if seen[n] {
			t.Fatalf("%s: repeating members: %s", in, q.String())
		}
		seen[n] = true
		if s[i] != n {
			t.Fatalf("slice has bad contents: %v != %v", s, q)
		}
	}
}

var testcases = []string{
	"+*/-",
	"++-",
	"90",
	"123456789--",
	"8++-++-++0",
	"0",
	"AB--/CDEF",
}

func TestDeque(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc, func(t *testing.T) {
			dequeLang(t, tc)
		})
	}
}

func FuzzDeque(f *testing.F) {
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(dequeLang)
}

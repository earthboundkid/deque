//go:build goexperiment.rangefunc

package deque_test

import (
	"container/list"
	"fmt"
	"math/rand/v2"
	"slices"
	"sort"
	"testing"

	"github.com/earthboundkid/deque/v2"
)

func ExampleDeque() {
	// Make a new deque
	d := deque.Of(9, 8, 7, 6)
	// Sort it
	sort.Sort(deque.Sortable[int]{d})
	// Add 5, 4, 3, 2, 1 to the front
	for i := 5; i > 0; i-- {
		d.PushFront(i)
	}
	fmt.Println(d)
	// Now reverse loop through items
	for _, n := range d.Reverse() {
		fmt.Print(n, " ")
	}
	fmt.Println()
	// Output:
	// Deque{ len: 9, cap: 16, items: [1, 2, 3, 4, 5, 6, 7, 8, 9]}
	// 9 8 7 6 5 4 3 2 1
}

func dequeLang(t *testing.T, in string) {
	q := deque.Make[int](0)
	l := list.New()
	qlen := 0
	mincap := 0
	i := 0
	for _, c := range in {
		switch c {
		case '+':
			l.PushFront(i)
			q.PushFront(i)
			qlen++
			if mincap < qlen {
				mincap++
			}
		case '*':
			l.PushBack(i)
			q.PushBack(i)
			qlen++
			if mincap < qlen {
				mincap++
			}
		case '-':
			if n := l.Front(); n != nil {
				l.Remove(n)
			}
			q.RemoveFront()
			qlen--
			if qlen < 0 {
				qlen = 0
			}
		case '/':
			if n := l.Back(); n != nil {
				l.Remove(n)
			}
			q.RemoveBack()
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
				l.PushBack(i)
				s = append(s, i)
				i++
			}
			q.PushBackSlice(s)
			qlen += n
			if mincap < qlen {
				mincap = qlen
			}

		}
		i++
	}
	if llen := l.Len(); q.Len() != llen {
		t.Errorf("%s bad len %d != %d", in, q.Len(), llen)
	}
	if q.Len() != qlen {
		t.Errorf("%s bad len %d != %d", in, q.Len(), qlen)
	}
	if q.Cap() < mincap {
		t.Errorf("%s: bad cap %d < %d", in, q.Cap(), mincap)
	}

	for cursor, n := 0, l.Front(); n != nil; n = n.Next() {
		if v, _ := q.At(cursor); v != n.Value.(int) {
			t.Errorf("deque.At(%d) == %d; want %d", cursor, v, n.Value)
		}
		cursor++
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

func TestDequeBasics(t *testing.T) {
	q := deque.Make[int](1)
	v, ok := q.Front()
	if v != 0 || ok {
		t.Errorf("empty deque.Head() got %v %v", v, ok)
	}
	v, ok = q.Back()
	if v != 0 || ok {
		t.Errorf("empty deque.Tail() got %v %v", v, ok)
	}

	v, ok = q.At(0)
	if v != 0 || ok {
		t.Errorf("empty deque.Tail() got %v %v", v, ok)
	}

	q.PushFront(1)

	v, ok = q.Front()
	if v != 1 || !ok {
		t.Errorf("deque{1}.Head() got %v %v", v, ok)
	}
	v, ok = q.Back()
	if v != 1 || !ok {
		t.Errorf("deque{1}.Tail() got %v %v", v, ok)
	}
	v, ok = q.At(0)
	if v != 1 || !ok {
		t.Errorf("empty deque.Tail() got %v %v", v, ok)
	}
}

func FuzzPropTest(f *testing.F) {
	f.Add(uint64(0), uint64(0))
	f.Fuzz(func(t *testing.T, seed1, seed2 uint64) {
		r := rand.New(rand.NewPCG(seed1, seed2))
		var d deque.Deque[int]
		l := list.New()
		for {
			f := r.Float64()
			switch {
			case f < .04:
				s1 := d.Slice()
				var s3 []int
				for n := l.Front(); n != nil; n = n.Next() {
					s3 = append(s3, n.Value.(int))
				}
				if !slices.Equal(s1, s3) {
					t.Fatal(s1, s3)
				}
				return
			case f < .28:
				n := r.IntN(100)
				d.PushFront(n)
				l.PushFront(n)
			case f < .52:
				n := r.IntN(100)
				d.PushBack(n)
				l.PushBack(n)
			case f < .76:
				v1, ok1 := d.Front()
				v2, ok3 := d.RemoveFront()
				if v1 != v2 || ok1 != ok3 {
					t.Fatal(v1, v2)
				}
				n := l.Front()
				if ok1 && n == nil || !ok1 && n != nil {
					t.Fatal(d, l, v2)
				}
				if n != nil {
					v3 := l.Remove(n).(int)
					if v3 != v2 {
						t.Fatal(d, l, v2, v3)
					}
				}
			default:
				v1, ok1 := d.Back()
				v2, ok2 := d.RemoveBack()
				if v1 != v2 {
					t.Fatal(v1, v2)
				}
				if v1 != v2 || ok1 != ok2 {
					t.Fatal(v1, v2)
				}
				n := l.Back()
				if ok1 && n == nil || !ok1 && n != nil {
					t.Fatal(d, l, v2)
				}
				if n != nil {
					v3 := l.Remove(n).(int)
					if v3 != v2 {
						t.Fatal(d, l, v2, v3)
					}
				}
			}
		}
	})
}

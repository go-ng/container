// This code is a modified version of the file sort/slice.go
// file from Go sorce code with license header:
//
// > Copyright 2017 The Go Authors. All rights reserved.
// > Use of this source code is governed by a BSD-style
// > license that can be found in the LICENSE file.
//
// Use of this modified code is allowed with respect to this license.
// The "LICENSE file" could be found in this directory
// with name "GOLANG-LICENSE"
//
// Any rights to the modifications themselves are waived and are
// available also under CC-0 1.0 license.
//                                             -- Dmitrii Okunev

package heap

import (
	"math/rand"
	"sort"
	"testing"
)

type myHeap struct {
	Slice sort.IntSlice
}

func (h *myHeap) Less(i, j int) bool {
	return h.Slice[i] < h.Slice[j]
}

func (h *myHeap) Swap(i, j int) {
	h.Slice[i], h.Slice[j] = h.Slice[j], h.Slice[i]
}

func (h *myHeap) Len() int {
	return len(h.Slice)
}

func (h *myHeap) Pop() (v any) {
	h.Slice, v = h.Slice[:h.Len()-1], h.Slice[h.Len()-1]
	return
}

func (h *myHeap) Push(v any) {
	h.Slice = append(h.Slice, v.(int))
}

func (h myHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.Slice[i], j1, h.Slice[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.Slice[i], j1, h.Slice[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	h := new(myHeap)
	for i := 20; i > 0; i-- {
		h.Push(0) // all elements are the same
	}
	Init(h.Slice)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(&h.Slice)
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := new(myHeap)
	for i := 20; i > 0; i-- {
		h.Push(i) // all elements are different
	}
	Init(h.Slice)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(&h.Slice)
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	h := new(myHeap)
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	Init(h.Slice)
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		Push(&h.Slice, i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := Pop(&h.Slice)
		if i < 20 {
			Push(&h.Slice, 20+i)
		}
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := new(myHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := Remove(&h.Slice, i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := new(myHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := Remove(&h.Slice, 0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := new(myHeap)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[Remove(&h.Slice, (h.Len()-1)/2)] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	var h myHeap
	h.Slice = make([]int, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Push(&h.Slice, 0) // all elements are the same
		}
		for h.Len() > 0 {
			Pop(&h.Slice)
		}
	}
}

func TestFix(t *testing.T) {
	var h myHeap
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		Push(&h.Slice, i)
	}
	h.verify(t, 0)

	if h.Slice[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", h.Slice[0])
	}
	h.Slice[0] = 210
	Fix(h.Slice, 0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			h.Slice[elem] *= 2
		} else {
			h.Slice[elem] /= 2
		}
		Fix(h.Slice, elem)
		h.verify(t, 0)
	}
}

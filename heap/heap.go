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

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//
package heap

type Interface[T any] interface {
	Less(i, j int) bool
	~[]T
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func Init[E any, S Interface[E]](s S) {
	// heapify
	n := len(s)
	for i := n/2 - 1; i >= 0; i-- {
		down(s, i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func Push[E any, S Interface[E]](s *S, x E) {
	*s = append(*s, x)
	up(*s, len(*s)-1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func Pop[E any, S Interface[E]](sP *S) E {
	s := *sP
	n := len(s) - 1
	s[0], s[n] = s[n], s[0]
	down(s, 0, n)

	v := s[n]
	*sP = s[0:n]
	return v
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func Remove[E any, S Interface[E]](s *S, i int) E {
	n := len(*s) - 1
	if n != i {
		(*s)[i], (*s)[n] = (*s)[n], (*s)[i]
		if !down(*s, i, n) {
			up(*s, i)
		}
	}

	v := (*s)[n]
	*s = (*s)[0:n]
	return v
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func Fix[E any, S Interface[E]](s S, i int) {
	if !down(s, i, len(s)) {
		up(s, i)
	}
}

func up[E any, S Interface[E]](s S, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !s.Less(j, i) {
			break
		}
		s[i], s[j] = s[j], s[i]
		j = i
	}
}

func down[E any, S Interface[E]](s S, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && s.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !s.Less(j, i) {
			break
		}
		s[i], s[j] = s[j], s[i]
		i = j
	}
	return i > i0
}

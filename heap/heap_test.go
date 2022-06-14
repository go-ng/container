// This file is available under CC-0 1.0 license.
//
// See file `CC0-LICENSE`.

package heap

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *intHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *intHeap) Pop() any {
	n := len(*h)
	v := (*h)[n-1]
	*h = (*h)[0 : n-1]
	return v
}

func BenchmarkStdHeap(b *testing.B) {
	for size := 1; size <= 1024*1024; size *= 2 {
		b.Run(fmt.Sprintf("size%d", size), func(b *testing.B) {
			for _, isPreallocated := range []bool{true, false} {
				b.Run(fmt.Sprintf("preallocated-%v", isPreallocated), func(b *testing.B) {
					var h intHeap
					if isPreallocated {
						h = make(intHeap, 0, size)
					}
					b.Run("Fill", func(b *testing.B) {
						rng := rand.New(rand.NewSource(0))
						values := rng.Perm(size)
						b.ReportAllocs()
						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							if isPreallocated {
								h = h[:0]
							} else {
								h = nil
							}
							for idx := 0; idx < size; idx++ {
								heap.Push(&h, values[idx])
							}
						}
					})
				})
			}
			b.Run("PopPush", func(b *testing.B) {
				rng := rand.New(rand.NewSource(0))
				v := rng.Perm(size)
				h := intHeap(rng.Perm(size))
				heap.Init(&h)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					heap.Pop(&h)
					heap.Push(&h, v[i%size])
				}
			})
		})
	}
}

func BenchmarkOurHeap(b *testing.B) {
	for size := 1; size <= 1024*1024; size *= 2 {
		b.Run(fmt.Sprintf("size%d", size), func(b *testing.B) {
			for _, isPreallocated := range []bool{true, false} {
				b.Run(fmt.Sprintf("preallocated-%v", isPreallocated), func(b *testing.B) {
					var h sort.IntSlice
					if isPreallocated {
						h = make(sort.IntSlice, 0, size)
					}
					b.Run("Fill", func(b *testing.B) {
						rng := rand.New(rand.NewSource(0))
						values := rng.Perm(size)
						b.ReportAllocs()
						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							if isPreallocated {
								h = h[:0]
							} else {
								h = nil
							}
							for idx := 0; idx < size; idx++ {
								Push(&h, values[idx])
							}
						}
					})
				})
			}
			b.Run("PopPush", func(b *testing.B) {
				rng := rand.New(rand.NewSource(0))
				v := rng.Perm(size)
				h := sort.IntSlice(rng.Perm(size))
				Init(h)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					Pop(&h)
					Push(&h, v[i%size])
				}
			})
		})
	}
}

type intSlice []int

func (s *intSlice) Less(i, j int) bool {
	return (*s)[i] < (*s)[j]
}

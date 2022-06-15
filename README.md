[![GoDoc](https://godoc.org/github.com/go-ng/container/heap?status.svg)](https://pkg.go.dev/github.com/go-ng/container/heap?tab=doc)

# About

This package just copies the original Go's package `container/heap`, but provides a generics-based API instead of an interface-argument based, which results in better performance.

An example of a migration. With standard `container/heap`:
```go
// This example demonstrates an integer heap built using the heap interface.
package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h).(int))
	}
}
```

A migrated code:

```go
// This example demonstrates an integer heap built using the heap interface.
package main

import (
	"fmt"

	"github.com/go-ng/container/heap"
)

type IntHeap []int

func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
	
func main() {
	h := IntHeap{2, 1, 5}
	heap.Init(&h)
	heap.Push(&h, 3)
	fmt.Printf("minimum: %d\n", h[0])
	for len(h) > 0 {
		fmt.Printf("%d ", heap.Pop(&h)) // Pop will return an `int`, not an `interface{}`
	}
}
```

That's it.

# Benchmarks

### `container/heap`

On my laptop it takes 15%-65% less CPU than the standard heap implementation (see `heap/BENCHMARKS.md`).

```
name                                         old time/op    new time/op    delta
Heap/size1/preallocated-true/Fill-16           7.37ns ± 4%    3.62ns ± 4%   -50.85%  (p=0.008 n=5+5)
Heap/size1/preallocated-false/Fill-16          23.3ns ±13%    19.1ns ± 3%   -18.02%  (p=0.008 n=5+5)
Heap/size1/PopPush-16                          14.2ns ± 3%     7.1ns ± 2%   -49.79%  (p=0.008 n=5+5)
Heap/size2/preallocated-true/Fill-16           17.4ns ± 4%    10.1ns ± 3%   -42.03%  (p=0.008 n=5+5)
Heap/size2/preallocated-false/Fill-16          55.5ns ± 4%    47.0ns ± 6%   -15.36%  (p=0.008 n=5+5)
Heap/size2/PopPush-16                          16.5ns ± 6%    10.4ns ± 7%   -37.30%  (p=0.008 n=5+5)
Heap/size4/preallocated-true/Fill-16           35.6ns ± 6%    23.5ns ± 4%   -34.08%  (p=0.008 n=5+5)
Heap/size4/preallocated-false/Fill-16          83.4ns ± 2%    70.9ns ± 4%   -14.98%  (p=0.008 n=5+5)
Heap/size4/PopPush-16                          23.4ns ± 3%    17.1ns ± 5%   -26.86%  (p=0.008 n=5+5)
Heap/size8/preallocated-true/Fill-16           73.8ns ± 3%    44.4ns ± 1%   -39.88%  (p=0.008 n=5+5)
Heap/size8/preallocated-false/Fill-16           165ns ± 2%     133ns ± 6%   -19.20%  (p=0.008 n=5+5)
Heap/size8/PopPush-16                          29.3ns ± 3%    20.9ns ± 4%   -28.69%  (p=0.008 n=5+5)
Heap/size16/preallocated-true/Fill-16           147ns ± 1%      94ns ± 1%   -36.45%  (p=0.008 n=5+5)
Heap/size16/preallocated-false/Fill-16          249ns ± 3%     199ns ± 6%   -20.11%  (p=0.008 n=5+5)
Heap/size16/PopPush-16                         31.0ns ± 1%    24.2ns ± 2%   -22.05%  (p=0.008 n=5+5)
Heap/size32/preallocated-true/Fill-16           376ns ± 6%     240ns ± 3%   -36.18%  (p=0.008 n=5+5)
Heap/size32/preallocated-false/Fill-16          522ns ± 5%     405ns ± 9%   -22.47%  (p=0.016 n=4+5)
Heap/size32/PopPush-16                         35.6ns ± 2%    25.4ns ± 4%   -28.76%  (p=0.008 n=5+5)
Heap/size64/preallocated-true/Fill-16           749ns ± 6%     473ns ± 4%   -36.89%  (p=0.008 n=5+5)
Heap/size64/preallocated-false/Fill-16          957ns ± 2%     677ns ± 1%   -29.26%  (p=0.016 n=5+4)
Heap/size64/PopPush-16                         39.0ns ± 0%    28.9ns ± 3%   -26.03%  (p=0.016 n=4+5)
Heap/size128/preallocated-true/Fill-16         1.47µs ± 3%    0.99µs ± 4%   -33.04%  (p=0.008 n=5+5)
Heap/size128/preallocated-false/Fill-16        1.81µs ± 2%    1.26µs ± 3%   -30.49%  (p=0.008 n=5+5)
Heap/size128/PopPush-16                        39.1ns ± 0%    31.3ns ± 3%   -20.02%  (p=0.016 n=4+5)
Heap/size256/preallocated-true/Fill-16         2.98µs ± 4%    1.95µs ± 3%   -34.68%  (p=0.008 n=5+5)
Heap/size256/preallocated-false/Fill-16        3.38µs ± 1%    2.41µs ± 4%   -28.48%  (p=0.008 n=5+5)
Heap/size256/PopPush-16                        45.0ns ± 4%    34.2ns ± 1%   -23.90%  (p=0.008 n=5+5)
Heap/size512/preallocated-true/Fill-16         8.24µs ± 1%    4.19µs ± 2%   -49.19%  (p=0.008 n=5+5)
Heap/size512/preallocated-false/Fill-16        11.3µs ± 7%     5.0µs ± 2%   -56.07%  (p=0.008 n=5+5)
Heap/size512/PopPush-16                        60.0ns ± 1%    35.8ns ± 2%   -40.34%  (p=0.008 n=5+5)
Heap/size1024/preallocated-true/Fill-16        22.6µs ± 1%     8.5µs ± 1%   -62.57%  (p=0.016 n=4+5)
Heap/size1024/preallocated-false/Fill-16       28.4µs ± 2%    11.4µs ± 2%   -59.81%  (p=0.008 n=5+5)
Heap/size1024/PopPush-16                       67.4ns ± 2%    41.2ns ± 2%   -38.92%  (p=0.008 n=5+5)
Heap/size2048/preallocated-true/Fill-16        52.9µs ± 4%    18.8µs ± 0%   -64.46%  (p=0.016 n=5+4)
Heap/size2048/preallocated-false/Fill-16       61.5µs ± 5%    27.4µs ± 1%   -55.49%  (p=0.008 n=5+5)
Heap/size2048/PopPush-16                       74.4ns ± 3%    45.2ns ± 3%   -39.31%  (p=0.008 n=5+5)
Heap/size4096/preallocated-true/Fill-16         112µs ± 4%      53µs ± 2%   -52.74%  (p=0.008 n=5+5)
Heap/size4096/preallocated-false/Fill-16        129µs ± 3%      65µs ± 1%   -49.28%  (p=0.008 n=5+5)
Heap/size4096/PopPush-16                       89.8ns ± 1%    59.7ns ± 2%   -33.49%  (p=0.008 n=5+5)
Heap/size8192/preallocated-true/Fill-16         233µs ± 3%     120µs ± 1%   -48.78%  (p=0.008 n=5+5)
Heap/size8192/preallocated-false/Fill-16        262µs ± 3%     141µs ± 1%   -45.94%  (p=0.016 n=5+4)
Heap/size8192/PopPush-16                        108ns ± 6%      73ns ± 2%   -32.37%  (p=0.008 n=5+5)
Heap/size16384/preallocated-true/Fill-16        483µs ± 2%     242µs ± 2%   -49.93%  (p=0.008 n=5+5)
Heap/size16384/preallocated-false/Fill-16       530µs ± 3%     289µs ± 3%   -45.51%  (p=0.008 n=5+5)
Heap/size16384/PopPush-16                       147ns ± 1%     111ns ± 1%   -24.32%  (p=0.008 n=5+5)
Heap/size32768/preallocated-true/Fill-16        975µs ± 4%     503µs ± 2%   -48.43%  (p=0.008 n=5+5)
Heap/size32768/preallocated-false/Fill-16      1.12ms ± 1%    0.60ms ± 1%   -46.42%  (p=0.016 n=5+4)
Heap/size32768/PopPush-16                       163ns ± 2%     129ns ± 4%   -21.11%  (p=0.008 n=5+5)
Heap/size65536/preallocated-true/Fill-16       1.94ms ± 1%    1.01ms ± 3%   -47.90%  (p=0.008 n=5+5)
Heap/size65536/preallocated-false/Fill-16      2.19ms ± 2%    1.24ms ± 1%   -43.29%  (p=0.008 n=5+5)
Heap/size65536/PopPush-16                       174ns ± 2%     139ns ± 3%   -20.17%  (p=0.008 n=5+5)
Heap/size131072/preallocated-true/Fill-16      3.90ms ± 3%    2.02ms ± 1%   -48.15%  (p=0.008 n=5+5)
Heap/size131072/preallocated-false/Fill-16     5.34ms ± 3%    2.76ms ± 1%   -48.39%  (p=0.008 n=5+5)
Heap/size131072/PopPush-16                      189ns ± 0%     152ns ± 3%   -20.00%  (p=0.008 n=5+5)
Heap/size262144/preallocated-true/Fill-16      7.80ms ± 3%    4.12ms ± 3%   -47.18%  (p=0.008 n=5+5)
Heap/size262144/preallocated-false/Fill-16     11.0ms ± 1%     5.1ms ± 1%   -53.16%  (p=0.008 n=5+5)
Heap/size262144/PopPush-16                      211ns ± 3%     165ns ± 2%   -21.81%  (p=0.008 n=5+5)
Heap/size524288/preallocated-true/Fill-16      15.9ms ± 3%     8.3ms ± 2%   -48.17%  (p=0.008 n=5+5)
Heap/size524288/preallocated-false/Fill-16     20.7ms ± 3%    11.4ms ±10%   -44.97%  (p=0.008 n=5+5)
Heap/size524288/PopPush-16                      232ns ± 2%     177ns ± 5%   -23.72%  (p=0.008 n=5+5)
Heap/size1048576/preallocated-true/Fill-16     31.5ms ± 2%    16.4ms ± 2%   -48.00%  (p=0.008 n=5+5)
Heap/size1048576/preallocated-false/Fill-16    39.1ms ± 3%    22.4ms ± 3%   -42.76%  (p=0.008 n=5+5)
Heap/size1048576/PopPush-16                     263ns ± 5%     199ns ± 5%   -24.37%  (p=0.008 n=5+5)
```
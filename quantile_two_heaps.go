package main

import "container/heap"

// Exact quantile by Two Heaps algorithm with linear interpolation
//
// Basic description could be found here:
// https://preparekaro.com/two-heap-approach.html
//
// Time: O(n log n) for adding n values and O(1) for computing the quantile.
// Mempory: O(n) for storing n values.
//
// NOT THREADSAFE
type QuantileTwoHeaps struct {
	left  MaxHeap
	right MinHeap
	p     float32
}

func NewQuantileTwoHeaps(p float32) QuantileTwoHeaps {
	return QuantileTwoHeaps{
		left:  make(MaxHeap, 0),
		right: make(MinHeap, 0),
		p:     p,
	}
}

// Push adds a new value to the quantile computation.
func (q *QuantileTwoHeaps) Push(val float32) {
	// Choosing target heap
	if q.left.Len() == 0 || val <= q.left[0] {
		heap.Push(&q.left, val)
	} else {
		heap.Push(&q.right, val)
	}

	totalN := q.left.Len() + q.right.Len()
	targetLeftSize := int(q.p*float32(totalN-1)) + 1

	// Rebalancing
	if q.left.Len() > targetLeftSize {
		v := heap.Pop(&q.left)
		heap.Push(&q.right, v)
	}
	if q.left.Len() < targetLeftSize {
		v := heap.Pop(&q.right)
		heap.Push(&q.left, v)
	}
}

// Value returns the current quantile value.
func (q *QuantileTwoHeaps) Value() float32 {
	// Edge cases
	if q.left.Len()+q.right.Len() == 0 {
		return 0
	}
	if q.right.Len() == 0 {
		return q.left[0]
	}
	if q.left.Len() == 0 {
		return q.right[0]
	}

	// Linear interpolation
	totalN := q.left.Len() + q.right.Len()
	pos := q.p * float32(totalN-1)
	lowerIdx := int(pos)
	fraction := pos - float32(lowerIdx)

	return q.left[0] + fraction*(q.right[0]-q.left[0])
}

// PushAndValue adds a new value and returns the current quantile value.
func (q *QuantileTwoHeaps) PushAndValue(value float32) float32 {
	q.Push(value)
	return q.Value()
}

// MinHeap is a min-heap implementation for container/heap.
type MinHeap []float32

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] } // Min
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(float32))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MaxHeap is a max-heap implementation for container/heap.
type MaxHeap []float32

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // Max
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(float32))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

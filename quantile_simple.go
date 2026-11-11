package main

import (
	"math"
	"slices"
)

// Exact quantile with straightforward implementaion with linear interpolation.
//
// Time: O(n log n) for adding n values and O(1) for computing the quantile.
// Mempory: O(n) for storing n values.
//
// It uses the following formulas:
//
//		  $$k = 1 + (N - 1) * p$$
//
//	where $p$ is the quantile to compute (0 <= p <= 1),
//
//		  $$q_p = x_i + d * (x_{i+1} - x_i)$$
//
// where $i$ - integer part of $k$, $d$ - rest part of $k$. $x_i$ is the $i$-th value in the sorted list of values.
//
// The implementation has low performance for large datasets, as it stores all values in memory and sorts them to
// compute the quantile. It is suitable for small datasets or when exact quantiles are required.
//
// NOT THREADSAFE
type QuantileSimple struct {
	p      float32
	values []float32
}

func NewQuantileSimple(p float32) QuantileSimple {
	return QuantileSimple{
		p: p,
	}
}

func (q *QuantileSimple) Push(val float32) {
	index, _ := slices.BinarySearch(q.values, val)
	q.values = slices.Insert(q.values, index, val)
}

func (q *QuantileSimple) Value() float32 {
	if len(q.values) == 0 {
		return 0
	}
	if q.p == 0 {
		return q.values[0]
	}
	if q.p == 1 {
		return q.values[len(q.values)-1]
	}

	k := 1 + (float32(len(q.values)-1))*q.p
	i := int64(math.Trunc(float64(k)))
	d := float32(k) - float32(i)

	xi := q.values[i-1]
	if d == 0 {
		return xi
	}

	return xi + d*(q.values[i]-xi)
}

func (q *QuantileSimple) PushAndValue(val float32) float32 {
	q.Push(val)
	return q.Value()
}

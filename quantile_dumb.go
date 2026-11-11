package main

// Dumb quantile implementation that does not store any values and always returns 0.
// This used for testing purposes where quantile computation is not required.
type QuantileDumb struct{}

func (q *QuantileDumb) Push(value float32) {}

func (q *QuantileDumb) Value(p float32) float32 {
	return 0
}

func (q *QuantileDumb) PushAndValue(value float32, p float32) float32 {
	return 0
}

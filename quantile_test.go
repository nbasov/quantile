package main

import (
	"math"
	"testing"
)

const quantileEps = 1e-5

// quantileImpls lists all quantile implementations that must satisfy the same test cases.
var quantileImpls = []struct {
	name string
	new  func(p float32) Computer
}{
	{
		name: "Simple",
		new: func(p float32) Computer {
			q := NewQuantileSimple(p)
			return &q
		},
	},
	{
		name: "TwoHeaps",
		new: func(p float32) Computer {
			q := NewQuantileTwoHeaps(p)
			return &q
		},
	},
}

func TestQuantile_PushAndValue(t *testing.T) {
	tests := []struct {
		name    string
		initial []float32
		value   float32
		p       float32
		want    float32
	}{
		{
			name:    "first value into empty",
			initial: nil,
			value:   5,
			p:       0.5,
			want:    5,
		},
		{
			name:    "median of even count interpolates",
			initial: []float32{1, 2, 3},
			value:   4,
			p:       0.5,
			want:    2.5,
		},
		{
			name:    "median of odd count is exact element",
			initial: []float32{3, 1},
			value:   2,
			p:       0.5,
			want:    2,
		},
		{
			name:    "p=0 returns new minimum",
			initial: []float32{5, 6},
			value:   1,
			p:       0,
			want:    1,
		},
		{
			name:    "p=1 returns new maximum",
			initial: []float32{5, 6},
			value:   10,
			p:       1,
			want:    10,
		},
		{
			name:    "value inserted in the middle keeps order",
			initial: []float32{1, 10},
			value:   5,
			p:       1,
			want:    10,
		},
		{
			name:    "duplicates",
			initial: []float32{1, 1},
			value:   1,
			p:       0.75,
			want:    1,
		},
		{
			name:    "negative values with interpolation",
			initial: []float32{-5, -1},
			value:   -3,
			p:       0.25,
			want:    -4,
		},
		{
			name:    "p90 of 1..10",
			initial: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
			value:   10,
			p:       0.9,
			want:    9.1,
		},
		{
			name:    "p25 of 1..5 hits exact element",
			initial: []float32{5, 4, 2, 1},
			value:   3,
			p:       0.25,
			want:    2,
		},
	}

	for _, impl := range quantileImpls {
		t.Run(impl.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					q := impl.new(tt.p)
					for _, v := range tt.initial {
						q.Push(v)
					}

					got := q.PushAndValue(tt.value)

					if math.Abs(float64(got-tt.want)) > quantileEps {
						t.Errorf("PushAndValue(%v) with p=%v = %v, want %v", tt.value, tt.p, got, tt.want)
					}
				})
			}
		})
	}
}

func TestQuantile_PushAndValue_Sequence(t *testing.T) {
	steps := []struct {
		value float32
		want  float32
	}{
		{value: 10, want: 10},  // [10]
		{value: 20, want: 15},  // [10 20]
		{value: 0, want: 10},   // [0 10 20]
		{value: 30, want: 15},  // [0 10 20 30]
		{value: -10, want: 10}, // [-10 0 10 20 30]
	}

	for _, impl := range quantileImpls {
		t.Run(impl.name, func(t *testing.T) {
			q := impl.new(0.5)
			for i, s := range steps {
				got := q.PushAndValue(s.value)
				if math.Abs(float64(got-s.want)) > quantileEps {
					t.Errorf("step %d: PushAndValue(%v) with p=0.5 = %v, want %v", i, s.value, got, s.want)
				}
			}
		})
	}
}

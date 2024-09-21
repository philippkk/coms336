package utils

import "math"

var Empty = Interval{Min: math.Inf(+1), Max: math.Inf(-1)}
var Universe = Interval{Min: math.Inf(-1), Max: math.Inf(1)}

type Interval struct {
	Min, Max float64
}

func (i Interval) size() float64 {
	return i.Max - i.Min
}
func (i Interval) contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}
func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

func (i Interval) clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

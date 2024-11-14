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

func NewIntervalFromIntervals(a, b Interval) Interval {
	var m, f float64

	if a.Min <= b.Min {
		m = a.Min
	} else {
		m = b.Min
	}

	if a.Max >= b.Max {
		f = a.Max
	} else {
		f = b.Max
	}

	return Interval{m, f}
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

func (i Interval) expand(delta float64) Interval {
	padding := delta / 2
	return Interval{i.Min - padding, i.Max + padding}
}

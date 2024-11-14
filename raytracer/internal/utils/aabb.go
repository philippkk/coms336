package utils

// AABB represents an Axis-Aligned Bounding Box
type AABB struct {
	X, Y, Z Interval
}

func NewAABB() AABB {
	return AABB{
		X: Interval{},
		Y: Interval{},
		Z: Interval{},
	}
}

func NewAABBFromIntervals(x, y, z Interval) AABB {
	return AABB{
		X: x,
		Y: y,
		Z: z,
	}
}

func NewAABBFromPoints(a, b Vec3) AABB {
	var x, y, z Interval

	if a.X <= b.X {
		x = Interval{a.X, b.X}
	} else {
		x = Interval{b.X, a.X}
	}
	if a.Y <= b.Y {
		y = Interval{a.Y, b.Y}
	} else {
		y = Interval{b.Y, a.Y}
	}
	if a.Z <= b.Z {
		z = Interval{a.Z, b.Z}
	} else {
		z = Interval{b.Z, a.Z}
	}
	return AABB{X: x, Y: y, Z: z}
}

// AxisInterval returns the interval for the specified axis
func (a *AABB) AxisInterval(n int) Interval {
	switch n {
	case 1:
		return a.Y
	case 2:
		return a.Z
	default:
		return a.X
	}
}

// Hit tests if a ray intersects with the AABB
func (a *AABB) Hit(r *Ray, rayT Interval) bool {
	rayOrig := r.Origin
	rayDir := r.Direction

	for axis := 0; axis < 3; axis++ {
		ax := a.AxisInterval(axis)
		var rayOrigAxis, rayDirAxis float64

		// Get the appropriate component based on axis
		switch axis {
		case 0:
			rayOrigAxis = rayOrig.X
			rayDirAxis = rayDir.X
		case 1:
			rayOrigAxis = rayOrig.Y
			rayDirAxis = rayDir.Y
		case 2:
			rayOrigAxis = rayOrig.Z
			rayDirAxis = rayDir.Z
		}

		adinv := 1.0 / rayDirAxis
		t0 := (ax.Min - rayOrigAxis) * adinv
		t1 := (ax.Max - rayOrigAxis) * adinv

		if t0 < t1 {
			if t0 > rayT.Min {
				rayT.Min = t0
			}
			if t1 < rayT.Max {
				rayT.Max = t1
			}
		} else {
			if t1 > rayT.Min {
				rayT.Min = t1
			}
			if t0 < rayT.Max {
				rayT.Max = t0
			}
		}
		if rayT.Max <= rayT.Min {
			return false
		}
	}

	return true
}

// Helper functions for creating combined AABBs

// SurroundingBox returns an AABB that contains both input boxes
func SurroundingBox(box0, box1 AABB) AABB {
	x := NewIntervalFromIntervals(box0.X, box1.X)

	y := NewIntervalFromIntervals(box0.Y, box1.Y)

	z := NewIntervalFromIntervals(box0.Z, box1.Z)

	return AABB{X: x, Y: y, Z: z}
}

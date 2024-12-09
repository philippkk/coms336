package utils

// AABB represents an Axis-Aligned Bounding Box
type AABB struct {
	X, Y, Z Interval
}

func (a *AABB) LongestAxis() int {
	if a.X.size() > a.Y.size() {
		if a.X.size() > a.Z.size() {
			return 0
		}
	} else {
		if a.Y.size() > a.Z.size() {
			return 1
		}
	}

	return 2
}
func (a *AABB) PadToMinimums() {
	delta := 0.0001
	if a.X.size() < delta {
		a.X = a.X.expand(delta)
	}
	if a.Y.size() < delta {
		a.Y = a.Y.expand(delta)
	}
	if a.Z.size() < delta {
		a.Z = a.Z.expand(delta)
	}
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
	thing := AABB{X: x, Y: y, Z: z}
	thing.PadToMinimums()
	return thing
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
	tMin := rayT.Min
	tMax := rayT.Max

	for axis := 0; axis < 3; axis++ {
		ax := a.AxisInterval(axis)
		rayOrigAxis := r.Origin.Get(axis)
		rayDirAxis := r.Direction.Get(axis)

		adinv := 1.0 / rayDirAxis
		t0 := (ax.Min - rayOrigAxis) * adinv
		t1 := (ax.Max - rayOrigAxis) * adinv

		if adinv < 0 {
			t0, t1 = t1, t0
		}

		tMin = max(t0, tMin)
		tMax = min(t1, tMax)

		if tMax <= tMin {
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

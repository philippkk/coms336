package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Triangle struct {
	V0, V1, V2 utils.Vec3 // Vertices of the triangle
	Mat        utils.Material
}

func (t Triangle) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	edge1 := t.V1.MinusEq(t.V0)
	edge2 := t.V2.MinusEq(t.V0)

	h := ray.Direction.Cross(edge2)
	a := edge1.Dot(h)

	epsilon := 1e-8
	if math.Abs(a) < epsilon {
		return false
	}

	f := 1.0 / a
	s := ray.Origin.MinusEq(t.V0)
	u := f * s.Dot(h)

	// Ray lies outside the triangle
	if u < 0.0 || u > 1.0 {
		return false
	}

	q := s.Cross(edge1)
	v := f * ray.Direction.Dot(q)

	// Ray lies outside the triangle
	if v < 0.0 || u+v > 1.0 {
		return false
	}

	// Calculate t value
	tval := f * edge2.Dot(q)

	if !rayT.Surrounds(tval) {
		return false
	}

	// Ray intersection
	rec.T = tval
	rec.P = ray.At(rec.T)

	// Calculate normal (ensure it's normalized)
	normal := edge1.Cross(edge2).Normalize()
	rec.SetFaceNormal(ray, normal)
	rec.Mat = t.Mat

	return true
}

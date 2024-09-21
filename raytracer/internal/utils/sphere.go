package utils

import "math"

type Sphere struct {
	Center Vec3
	Radius float64
}

func (s Sphere) Hit(ray *Ray, rayTmin, rayTmax float64, rec HitRecord) bool {
	oc := s.Center.MinusEq(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant > 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	root := (h - sqrtd) / a
	if root <= rayTmin || rayTmax <= root {
		root = (h + sqrtd) / a
		if root <= rayTmin || rayTmax <= root {
			return false
		}
	}

	rec.T = root
	rec.P = ray.At(rec.T)
	outwardNormal := (rec.P.MinusEq(s.Center)).TimesConst(1.0 / s.Radius)
	rec.SetFaceNormal(ray, outwardNormal)

	return true
}

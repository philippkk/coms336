package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Sphere struct {
	Center utils.Vec3
	Radius float64
	Mat    utils.Material
}

func (s Sphere) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	oc := s.Center.MinusEq(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = ray.At(rec.T)
	outwardNormal := (rec.P.MinusEq(s.Center)).TimesConst(1.0 / s.Radius)
	rec.SetFaceNormal(ray, outwardNormal)
	rec.Mat = s.Mat

	return true
}

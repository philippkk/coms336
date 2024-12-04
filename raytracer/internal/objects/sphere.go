package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

// Moving sphere is center1, center2 - center1, 0 time
type Sphere struct {
	Center utils.Ray
	Radius float64
	Mat    utils.Material
	Bbox   utils.AABB
}

func CreateSphere(Center utils.Ray, Radius float64, Mat utils.Material) Sphere {
	rvec := utils.Vec3{X: Radius, Y: Radius, Z: Radius}
	box1 := utils.NewAABBFromPoints(Center.At(0).MinusEq(rvec), Center.At(0).PlusEq(rvec))
	box2 := utils.NewAABBFromPoints(Center.At(1).MinusEq(rvec), Center.At(1).PlusEq(rvec))
	bbox := utils.SurroundingBox(box1, box2)
	return Sphere{Center, Radius, Mat, bbox}
}
func (s Sphere) BoundingBox() utils.AABB {
	return s.Bbox
}
func (s Sphere) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	currentCenter := s.Center.At(ray.Tm)
	oc := currentCenter.MinusEq(ray.Origin)
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
	outwardNormal := (rec.P.MinusEq(currentCenter)).TimesConst(1.0 / s.Radius)
	rec.SetFaceNormal(ray, outwardNormal)
	rec.U, rec.V = SphereUV(outwardNormal)
	rec.Mat = s.Mat

	return true
}

func SphereUV(p utils.Vec3) (u, v float64) {
	theta := math.Acos(p.Y / p.Length())
	phi := math.Atan2(p.Z, p.X)
	u = 0.5 + phi/(2*math.Pi)
	v = 1 - theta/math.Pi
	return
}

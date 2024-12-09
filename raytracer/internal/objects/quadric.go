package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type QuadricType int

const (
	SPHERE QuadricType = iota
	CYLINDER
	CONE
	ELLIPSOID
	HYPERBOLOID
)

type Quadric struct {
	A, B, C float64
	D, E, F float64
	G, H, I float64
	J       float64
	Mat     utils.Material
	Center  utils.Vec3
	Type    QuadricType
}

func CreateQuadricSphere(center utils.Vec3, radius float64, mat utils.Material) Quadric {
	return Quadric{
		A: 1, B: 1, C: 1,
		D: 0, E: 0, F: 0,
		G:      -2 * center.X,
		H:      -2 * center.Y,
		I:      -2 * center.Z,
		J:      center.X*center.X + center.Y*center.Y + center.Z*center.Z - radius*radius,
		Mat:    mat,
		Center: center,
		Type:   SPHERE,
	}
}

// CreateCylinder creates an infinite cylinder along the y-axis
func CreateCylinder(center utils.Vec3, radius float64, mat utils.Material) Quadric {
	return Quadric{
		A: 1, B: 0, C: 1, // x² + z² coefficients (no y² for infinite cylinder)
		D: 0, E: 0, F: 0, // no cross terms
		G:      -2 * center.X, // -2x₀x
		H:      0,             // no y term
		I:      -2 * center.Z, // -2z₀z
		J:      center.X*center.X + center.Z*center.Z - radius*radius,
		Mat:    mat,
		Center: center,
		Type:   CYLINDER,
	}
}

func CreateCone(center utils.Vec3, angle float64, mat utils.Material) Quadric {
	k := math.Tan(angle) * math.Tan(angle)
	return Quadric{
		A: 1, B: -k, C: 1,
		D: 0, E: 0, F: 0,
		G:      -2 * center.X,
		H:      2 * k * center.Y,
		I:      -2 * center.Z,
		J:      center.X*center.X - k*center.Y*center.Y + center.Z*center.Z,
		Mat:    mat,
		Center: center,
		Type:   CONE,
	}
}

func (q Quadric) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	origin := ray.Origin.MinusEq(q.Center)
	dir := ray.Direction

	a := q.A*dir.X*dir.X + q.B*dir.Y*dir.Y + q.C*dir.Z*dir.Z +
		q.D*dir.X*dir.Y + q.E*dir.X*dir.Z + q.F*dir.Y*dir.Z

	b := 2*q.A*origin.X*dir.X + 2*q.B*origin.Y*dir.Y + 2*q.C*origin.Z*dir.Z +
		q.D*(origin.X*dir.Y+origin.Y*dir.X) +
		q.E*(origin.X*dir.Z+origin.Z*dir.X) +
		q.F*(origin.Y*dir.Z+origin.Z*dir.Y) +
		q.G*dir.X + q.H*dir.Y + q.I*dir.Z

	c := q.A*origin.X*origin.X + q.B*origin.Y*origin.Y + q.C*origin.Z*origin.Z +
		q.D*origin.X*origin.Y + q.E*origin.X*origin.Z + q.F*origin.Y*origin.Z +
		q.G*origin.X + q.H*origin.Y + q.I*origin.Z + q.J

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)
	root := (-b - sqrtd) / (2 * a)

	if !rayT.Surrounds(root) {
		root = (-b + sqrtd) / (2 * a)
		if !rayT.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = ray.At(root)

	// Calculate normal
	normal := utils.Vec3{
		X: 2*q.A*rec.P.X + q.D*rec.P.Y + q.E*rec.P.Z + q.G,
		Y: 2*q.B*rec.P.Y + q.D*rec.P.X + q.F*rec.P.Z + q.H,
		Z: 2*q.C*rec.P.Z + q.E*rec.P.X + q.F*rec.P.Y + q.I,
	}

	rec.SetFaceNormal(ray, normal.Normalize())
	rec.Mat = q.Mat

	return true
}

func (q Quadric) BoundingBox() utils.AABB {
	bound := 100.0
	return utils.NewAABBFromPoints(
		utils.Vec3{X: q.Center.X - bound, Y: q.Center.Y - bound, Z: q.Center.Z - bound},
		utils.Vec3{X: q.Center.X + bound, Y: q.Center.Y + bound, Z: q.Center.Z + bound},
	)
}

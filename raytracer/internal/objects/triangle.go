package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Triangle struct {
	v0, v1, v2         utils.Vec3
	n0, n1, n2         utils.Vec3 // vertex normals
	uv0, uv1, uv2      utils.Vec2
	Mat                utils.Material
	Bbox               utils.AABB
	useSmoothedNormals bool
}

func (t Triangle) BoundingBox() utils.AABB {
	return t.Bbox
}
func CreateTriangleWithUV(v0, v1, v2 utils.Vec3, mat utils.Material, uv0, uv1, uv2 utils.Vec2) Triangle {
	// Existing bounding box creation
	min := utils.Vec3{
		X: math.Min(v0.X, math.Min(v1.X, v2.X)),
		Y: math.Min(v0.Y, math.Min(v1.Y, v2.Y)),
		Z: math.Min(v0.Z, math.Min(v1.Z, v2.Z)),
	}
	max := utils.Vec3{
		X: math.Max(v0.X, math.Max(v1.X, v2.X)),
		Y: math.Max(v0.Y, math.Max(v1.Y, v2.Y)),
		Z: math.Max(v0.Z, math.Max(v1.Z, v2.Z)),
	}

	bbox := utils.NewAABBFromPoints(min, max)

	return Triangle{
		v0: v0, v1: v1, v2: v2,
		uv0: uv0, uv1: uv1, uv2: uv2,
		Mat:  mat,
		Bbox: bbox,
	}
}
func CreateTriangleWithNormals(v0, v1, v2 utils.Vec3, n0, n1, n2 utils.Vec3, mat utils.Material, uv0, uv1, uv2 utils.Vec2) Triangle {
	min := utils.Vec3{
		X: math.Min(v0.X, math.Min(v1.X, v2.X)),
		Y: math.Min(v0.Y, math.Min(v1.Y, v2.Y)),
		Z: math.Min(v0.Z, math.Min(v1.Z, v2.Z)),
	}
	max := utils.Vec3{
		X: math.Max(v0.X, math.Max(v1.X, v2.X)),
		Y: math.Max(v0.Y, math.Max(v1.Y, v2.Y)),
		Z: math.Max(v0.Z, math.Max(v1.Z, v2.Z)),
	}

	bbox := utils.NewAABBFromPoints(min, max)

	return Triangle{
		v0: v0, v1: v1, v2: v2,
		n0: n0.Normalize(), n1: n1.Normalize(), n2: n2.Normalize(),
		uv0: uv0, uv1: uv1, uv2: uv2,
		Mat:                mat,
		Bbox:               bbox,
		useSmoothedNormals: true,
	}
}
func CreateTriangle(v0, v1, v2 utils.Vec3, mat utils.Material) Triangle {
	min := utils.Vec3{
		X: math.Min(v0.X, math.Min(v1.X, v2.X)),
		Y: math.Min(v0.Y, math.Min(v1.Y, v2.Y)),
		Z: math.Min(v0.Z, math.Min(v1.Z, v2.Z)),
	}
	max := utils.Vec3{
		X: math.Max(v0.X, math.Max(v1.X, v2.X)),
		Y: math.Max(v0.Y, math.Max(v1.Y, v2.Y)),
		Z: math.Max(v0.Z, math.Max(v1.Z, v2.Z)),
	}

	bbox := utils.NewAABBFromPoints(min, max)

	return Triangle{v0: v0, v1: v1, v2: v2, Mat: mat, Bbox: bbox}
}

func (t Triangle) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	edge1 := t.v1.MinusEq(t.v0)
	edge2 := t.v2.MinusEq(t.v0)

	h := ray.Direction.Cross(edge2)
	a := edge1.Dot(h)

	epsilon := 1e-8
	if math.Abs(a) < epsilon {
		return false
	}

	f := 1.0 / a
	s := ray.Origin.MinusEq(t.v0)
	u := f * s.Dot(h)

	if u < 0.0 || u > 1.0 {
		return false
	}

	q := s.Cross(edge1)
	v := f * ray.Direction.Dot(q)

	if v < 0.0 || u+v > 1.0 {
		return false
	}

	tval := f * edge2.Dot(q)

	if !rayT.Surrounds(tval) {
		return false
	}

	// Ray intersection
	rec.T = tval
	rec.P = ray.At(rec.T)

	// Calculate normal based on smoothing preference
	var normal utils.Vec3
	if t.useSmoothedNormals {
		// Interpolate normal using barycentric coordinates
		w := 1.0 - u - v
		normal = t.n0.TimesConst(w).
			PlusEq(t.n1.TimesConst(u)).
			PlusEq(t.n2.TimesConst(v)).
			Normalize()
	} else {
		// Use flat normal
		normal = edge1.Cross(edge2).Normalize()
	}

	rec.SetFaceNormal(ray, normal)
	rec.Mat = t.Mat

	// Interpolate UV coordinates
	w := 1.0 - u - v
	rec.U = w*t.uv0.X + u*t.uv1.X + v*t.uv2.X
	rec.V = w*t.uv0.Y + u*t.uv1.Y + v*t.uv2.Y

	return true
}

// Helper method to toggle between smooth and flat shading
func (t *Triangle) SetSmoothShading(enabled bool) {
	t.useSmoothedNormals = enabled
}

func mod1(x float64) float64 {
	// Handle both positive and negative UV coordinates
	x = math.Mod(x, 1.0)
	if x < 0 {
		x += 1.0
	}
	return x
}

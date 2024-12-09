package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Quad struct {
	Q      utils.Vec3 // Corner point
	U, V   utils.Vec3 // Two edges defining the quad
	Mat    utils.Material
	Normal utils.Vec3 // Normal vector
	D      float64    // Optional distance from origin
	W      utils.Vec3 // Normal basis vector
	Bbox   utils.AABB
}

func CreateQuad(q, u, v utils.Vec3, mat utils.Material) Quad {
	quad := Quad{
		Q:   q,
		U:   u,
		V:   v,
		Mat: mat,
	}

	// Calculate the normal vector from the edges
	n := u.Cross(v)

	// Set up the normal basis vectors
	quad.Normal = n.Normalize()
	quad.W = n.TimesConst(1.0 / n.Dot(n))

	// Set distance from origin
	quad.D = quad.Normal.Dot(q)

	// Calculate bounding box
	corner1 := q
	corner2 := q.PlusEq(u)
	corner3 := q.PlusEq(v)
	corner4 := q.PlusEq(u).PlusEq(v)

	min := utils.Vec3{
		X: math.Min(math.Min(corner1.X, corner2.X), math.Min(corner3.X, corner4.X)),
		Y: math.Min(math.Min(corner1.Y, corner2.Y), math.Min(corner3.Y, corner4.Y)),
		Z: math.Min(math.Min(corner1.Z, corner2.Z), math.Min(corner3.Z, corner4.Z)) - 0.0001, // Add small padding
	}
	max := utils.Vec3{
		X: math.Max(math.Max(corner1.X, corner2.X), math.Max(corner3.X, corner4.X)),
		Y: math.Max(math.Max(corner1.Y, corner2.Y), math.Max(corner3.Y, corner4.Y)),
		Z: math.Max(math.Max(corner1.Z, corner2.Z), math.Max(corner3.Z, corner4.Z)) + 0.0001, // Add small padding
	}

	quad.Bbox = utils.NewAABBFromPoints(min, max)
	return quad
}

func (q Quad) Hit(ray *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	denom := q.Normal.Dot(ray.Direction)

	// No hit if the ray is parallel to the plane
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// Calculate the distance to the plane
	t := (q.D - q.Normal.Dot(ray.Origin)) / denom

	// Return false if the intersection is outside the ray interval
	if !rayT.Surrounds(t) {
		return false
	}

	// Get the intersection point on the plane
	intersection := ray.At(t)

	// Get the hit point relative to the quad's corner
	planarPoint := intersection.MinusEq(q.Q)

	// Here's the key change: use the quadratic forms directly
	uLength := q.U.Length()
	vLength := q.V.Length()

	// Project onto U and V and normalize by their lengths
	u := planarPoint.Dot(q.U) / (uLength * uLength)
	v := planarPoint.Dot(q.V) / (vLength * vLength)

	// Check if we're inside the quad
	if u < 0 || u > 1 || v < 0 || v > 1 {
		return false
	}

	rec.T = t
	rec.P = intersection
	rec.Mat = q.Mat
	rec.U = u
	rec.V = v
	rec.SetFaceNormal(ray, q.Normal)

	return true
}

func (q Quad) BoundingBox() utils.AABB {
	return q.Bbox
}

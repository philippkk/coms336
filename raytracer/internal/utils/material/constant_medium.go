package material

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type ConstantMedium struct {
	Boundary      utils.Hittable
	NegInvDensity float64
	PhaseFunction utils.Material
}

func CreateConstantMedium(boundary utils.Hittable, density float64, albedo utils.Vec3) *ConstantMedium {
	return &ConstantMedium{
		Boundary:      boundary,
		NegInvDensity: -1.0 / density,
		PhaseFunction: NewIsotropicFromColor(albedo),
	}
}

func (cm *ConstantMedium) Hit(r *utils.Ray, rayT utils.Interval, rec *utils.HitRecord) bool {
	var rec1, rec2 utils.HitRecord

	universe := utils.Interval{Min: -math.Inf(1), Max: math.Inf(1)}
	if !cm.Boundary.Hit(r, universe, &rec1) {
		return false
	}

	if !cm.Boundary.Hit(r, utils.Interval{Min: rec1.T + 0.0001, Max: math.Inf(1)}, &rec2) {
		return false
	}

	if rec1.T < rayT.Min {
		rec1.T = rayT.Min
	}
	if rec2.T > rayT.Max {
		rec2.T = rayT.Max
	}

	if rec1.T >= rec2.T {
		return false
	}

	if rec1.T < 0 {
		rec1.T = 0
	}

	rayLength := r.Direction.Length()
	distanceInsideBoundary := (rec2.T - rec1.T) * rayLength
	hitDistance := cm.NegInvDensity * math.Log(utils.RandomFloat())

	if hitDistance > distanceInsideBoundary {
		return false
	}

	rec.T = rec1.T + hitDistance/rayLength
	rec.P = r.At(rec.T)

	rec.Normal = utils.Vec3{X: 1, Y: 0, Z: 0}
	rec.FrontFace = true
	rec.Mat = cm.PhaseFunction

	return true
}

func (cm *ConstantMedium) BoundingBox() utils.AABB {
	return cm.Boundary.BoundingBox()
}

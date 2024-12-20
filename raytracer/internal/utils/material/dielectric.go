package material

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Dielectric struct {
	RefractionIndex float64
}

func (d Dielectric) ColorEmitted(u, v float64, p utils.Vec3) utils.Vec3 {
	//TODO implement me
	return utils.Vec3{0, 0, 0}

}

func (d Dielectric) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	*attenuation = utils.Vec3{X: 1, Y: 1, Z: 1}
	var ri float64
	if rec.FrontFace {
		ri = 1.0 / d.RefractionIndex
	} else {
		ri = d.RefractionIndex
	}

	unitDirection := rIn.Direction.UnitVector()

	cosTheta := math.Min((unitDirection.Neg()).Dot(rec.Normal), 1)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1.0
	var direction utils.Vec3

	if cannotRefract || reflectance(cosTheta, ri) > utils.RandomFloat() {
		direction = utils.Reflect(unitDirection, rec.Normal)
	} else {
		direction = utils.Refract(unitDirection, rec.Normal, ri)
	}

	*scattered = utils.Ray{Origin: rec.P, Direction: direction, Tm: rIn.Tm}
	return true
}

func reflectance(cosine, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

package materials

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

type Dielectric struct {
	RefractionIndex float64
}

func (d Dielectric) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec utils.HitRecord) bool {
	*attenuation = utils.Vec3{1, 1, 1}
	var ri float64
	if rec.FrontFace {
		ri = 1.0 / d.RefractionIndex
	} else {
		ri = d.RefractionIndex
	}

	unitDirection := rIn.Direction.UnitVector()

	cosTheta := min(unitDirection.Neg().Dot(rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := ri*sinTheta > 1.0
	var direction utils.Vec3

	if cannotRefract {
		direction = utils.Reflect(unitDirection, rec.Normal)
	} else {
		direction = utils.Refract(unitDirection, rec.Normal, ri)
	}

	*scattered = utils.Ray{rec.P, direction}
	return true
}

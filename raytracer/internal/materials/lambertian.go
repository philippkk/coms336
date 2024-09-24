package materials

import "github.com/philippkk/coms336/raytracer/internal/utils"

type Lambertian struct {
	Albedo utils.Vec3
}

func (l Lambertian) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	scatterDirection := rec.Normal.PlusEq(utils.RandomUnitVector())

	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	*scattered = utils.Ray{Origin: rec.P, Direction: scatterDirection}
	*attenuation = l.Albedo
	return true
}

package materials

import "github.com/philippkk/coms336/raytracer/internal/utils"

type Metal struct {
	Albedo utils.Vec3
	Fuzz   float64
}

func (m Metal) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	reflected := utils.Reflect(rIn.Direction, rec.Normal)
	reflected = reflected.UnitVector().PlusEq(utils.RandomUnitVector().TimesConst(m.Fuzz))
	*scattered = utils.Ray{Origin: rec.P, Direction: reflected}
	*attenuation = m.Albedo
	return scattered.Direction.Dot(rec.Normal) > 0
}

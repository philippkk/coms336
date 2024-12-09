package material

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
)

type Lambertian struct {
	Tex utils.Texture
}

func (l Lambertian) ColorEmitted(u, v float64, p utils.Vec3) utils.Vec3 {
	return utils.Vec3{0, 0, 0}
}

func NewLambertianFromColor(albedo utils.Vec3) *Lambertian {
	return &Lambertian{Tex: utils.NewSolidColor(albedo)}
}

func NewLambertian(tex utils.Texture) *Lambertian {
	return &Lambertian{Tex: tex}
}

func (l Lambertian) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	scatterDirection := rec.Normal.PlusEq(utils.RandomUnitVector())

	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	*scattered = utils.Ray{Origin: rec.P, Direction: scatterDirection, Tm: rIn.Tm}

	*attenuation = l.Tex.Value(rec.U, rec.V, rec.P)
	return true
}

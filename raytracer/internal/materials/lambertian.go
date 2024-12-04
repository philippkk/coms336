package materials

import "github.com/philippkk/coms336/raytracer/internal/utils"

type Lambertian struct {
	Tex Texture
}

func NewLambertianFromColor(albedo utils.Vec3) *Lambertian {
	return &Lambertian{Tex: NewSolidColor(albedo)}
}

func NewLambertian(tex Texture) *Lambertian {
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

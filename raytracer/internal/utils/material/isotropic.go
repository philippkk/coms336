package material

import "github.com/philippkk/coms336/raytracer/internal/utils"

type Isotropic struct {
	texture utils.Texture
}

func (i Isotropic) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	scattered = &utils.Ray{Origin: rec.P, Direction: utils.RandomUnitVector(), Tm: rIn.Tm}
	*attenuation = i.texture.Value(rec.U, rec.V, rec.P)
	return true
}

func (i Isotropic) ColorEmitted(u, v float64, p utils.Vec3) utils.Vec3 {
	return utils.Vec3{}
}

func NewIsotropicFromColor(color utils.Vec3) *Isotropic {
	return &Isotropic{utils.NewSolidColor(color)}
}

func NewIsotropicFromTexture(texture utils.Texture) *Isotropic {
	return &Isotropic{texture}
}

package material

import "github.com/philippkk/coms336/raytracer/internal/utils"

type DiffuseLight struct {
	texture utils.Texture
}

func (d *DiffuseLight) Scatter(rIn, scattered *utils.Ray, attenuation *utils.Vec3, rec *utils.HitRecord) bool {
	return false
}

func NewDiffuseLightFromColor(color utils.Vec3) *DiffuseLight {
	return &DiffuseLight{utils.NewSolidColor(color)}
}

func NewDiffuseLightFromTexture(texture utils.Texture) *DiffuseLight {
	return &DiffuseLight{texture}
}

func (d *DiffuseLight) ColorEmitted(u, v float64, p utils.Vec3) utils.Vec3 {
	return d.texture.Value(u, v, p)
}

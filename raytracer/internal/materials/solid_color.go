package materials

import "github.com/philippkk/coms336/raytracer/internal/utils"

type SolidColor struct {
	Albedo utils.Vec3
}

func NewSolidColor(albedo utils.Vec3) *SolidColor {
	return &SolidColor{Albedo: albedo}
}

func NewSolidColorRGB(red, green, blue float64) *SolidColor {
	return &SolidColor{Albedo: utils.Vec3{red, green, blue}}
}

func (sc *SolidColor) Value(u, v float64, p utils.Vec3) utils.Vec3 {
	return sc.Albedo
}

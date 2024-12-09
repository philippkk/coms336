package utils

type SolidColor struct {
	Albedo Vec3
}

func NewSolidColor(albedo Vec3) *SolidColor {
	return &SolidColor{Albedo: albedo}
}

func NewSolidColorRGB(red, green, blue float64) *SolidColor {
	return &SolidColor{Albedo: Vec3{red, green, blue}}
}

func (sc *SolidColor) Value(u, v float64, p Vec3) Vec3 {
	return sc.Albedo
}

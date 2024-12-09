package utils

type Texture interface {
	Value(u, v float64, p Vec3) Vec3
}

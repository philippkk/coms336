package utils

type Vec3 struct {
	E [3]float64
}

func NewVec3(x, y, z float64) *Vec3 {
	var val [3]float64
	val[0] = x
	val[1] = y
	val[2] = z
	return &Vec3{val}
}

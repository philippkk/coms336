package vec3

type Vec3 struct {
	e [3]float64
}

func (v Vec3) NewVec3(x, y, z float64) *Vec3 {
	var val [3]float64
	val[0] = x
	val[1] = y
	val[2] = z
	return &Vec3{val}
}
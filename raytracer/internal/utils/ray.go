package utils

type Ray struct {
	Origin    Vec3
	Direction Vec3
	Tm        float64
}

func (r *Ray) At(t float64) Vec3 {
	return r.Origin.PlusEq(r.Direction.TimesConst(t))
}

package utils

import "math"

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
func (v *Vec3) Neg() *Vec3 {
	return NewVec3(-v.X(), -v.Y(), -v.Z())
}
func (v *Vec3) X() float64 { return v.E[0] }
func (v *Vec3) Y() float64 { return v.E[1] }
func (v *Vec3) Z() float64 { return v.E[2] }
func (v *Vec3) PlusEq(v2 *Vec3) {
	v.E[0] += v2.E[0]
	v.E[1] += v2.E[1]
	v.E[2] += v2.E[2]
}
func (v *Vec3) TimesEq(v2 *Vec3) {
	v.E[0] *= v2.E[0]
	v.E[1] *= v2.E[1]
	v.E[2] *= v2.E[2]
}
func (v *Vec3) TimesConst(t float64) {
	v.E[0] *= t
	v.E[1] *= t
	v.E[2] *= t
}
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}
func (v *Vec3) LengthSquared() float64 {
	return (v.X() * v.X()) +
		(v.Y() * v.Y()) +
		(v.Z() * v.Z())
}
func (v *Vec3) Dot(v2 *Vec3) float64 {
	return v.X()*v2.X() +
		v.Y()*v2.Y() +
		v.Z()*v2.Z()
}
func (v *Vec3) Cross(v2 *Vec3) *Vec3 {
	x := v.Y()*v2.Z() - v.Z()*v2.Y()
	y := v.Z()*v2.X() - v.X()*v2.Z()
	z := v.X()*v2.Y() - v.Y()*v2.X()
	return NewVec3(x, y, z)
}
func (v *Vec3) UnitVector() *Vec3 {
	v.TimesConst(1 / v.Length())
	return v
}

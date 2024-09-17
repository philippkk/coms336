package utils

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}
func (v Vec3) PlusEq(v2 Vec3) Vec3 {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}
func (v Vec3) PlusConst(v2 float64) Vec3 {
	v.X += v2
	v.Y += v2
	v.Z += v2
	return v
}
func (v Vec3) MinusEq(v2 Vec3) Vec3 {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}
func (v Vec3) TimesEq(v2 Vec3) Vec3 {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z
	return v
}
func (v Vec3) TimesConst(t float64) Vec3 {
	v.X *= t
	v.Y *= t
	v.Z *= t
	return v
}
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}
func (v Vec3) LengthSquared() float64 {
	return (v.X * v.X) +
		(v.Y * v.Y) +
		(v.Z * v.Z)
}
func (v Vec3) Dot(v2 Vec3) float64 {
	return v.X*v2.X +
		v.Y*v2.Y +
		v.Z*v2.Z
}
func (v Vec3) Cross(v2 Vec3) Vec3 {
	x := v.Y*v2.Z - v.Z*v2.Y
	y := v.Z*v2.X - v.X*v2.Z
	z := v.X*v2.Y - v.Y*v2.X
	return Vec3{x, y, z}
}
func (v Vec3) UnitVector() Vec3 {
	v.TimesConst(1.0 / v.Length())
	return v
}

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}

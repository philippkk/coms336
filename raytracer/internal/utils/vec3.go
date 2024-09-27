package utils

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}
func (v Vec3) PlusEq(o Vec3) Vec3 {
	return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}
func (v Vec3) PlusConst(t float64) Vec3 {
	return Vec3{v.X + t, v.Y + t, v.Z + t}
}
func (v Vec3) MinusEq(o Vec3) Vec3 {
	return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}
func (v Vec3) TimesEq(o Vec3) Vec3 {
	return Vec3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}
func (v Vec3) TimesConst(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}
func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}
func (v Vec3) LengthSquared() float64 {
	return (v.X * v.X) +
		(v.Y * v.Y) +
		(v.Z * v.Z)
}
func (v Vec3) Dot(o Vec3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}
func (v Vec3) Cross(v2 Vec3) Vec3 {
	x := v.Y*v2.Z - v.Z*v2.Y
	y := v.Z*v2.X - v.X*v2.Z
	z := v.X*v2.Y - v.Y*v2.X
	return Vec3{x, y, z}
}
func (v Vec3) UnitVector() Vec3 {
	return v.TimesConst(1.0 / v.Length())
}
func RandomInUnitDisk() Vec3 {
	for {
		p := Vec3{X: RandomFloatInRange(-1, 1), Y: RandomFloatInRange(-1, 1)}
		if p.LengthSquared() < 1 {
			return p
		}
	}
}
func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}
func (v Vec3) NearZero() bool {
	s := 1e-8
	return (math.Abs(v.X) < s) && (math.Abs(v.Y) < s) && (math.Abs(v.Z) < s)
}
func RandomVec3() Vec3 {
	return Vec3{RandomFloat(), RandomFloat(), RandomFloat()}
}
func RandomVec3InRange(min, max float64) Vec3 {
	return Vec3{RandomFloatInRange(min, max), RandomFloatInRange(min, max), RandomFloatInRange(min, max)}
}
func RandomUnitVector() Vec3 {
	for {
		p := RandomVec3InRange(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1 {
			return p.TimesConst(1 / math.Sqrt(lensq))
		}
	}
}
func RandomOnHemisphere(normal Vec3) Vec3 {
	onUnitSphere := RandomUnitVector()
	if onUnitSphere.Dot(normal) > 0.0 {
		return onUnitSphere
	} else {
		return onUnitSphere.Neg()
	}
}
func Reflect(v, n Vec3) Vec3 {
	return v.MinusEq(n.TimesConst(v.Dot(n) * 2))
}
func Refract(uv, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := min(n.Dot(uv.Neg()), 1.0)
	rOutPerp := uv.PlusEq(n.TimesConst(cosTheta)).TimesConst(etaiOverEtat)
	rOutParallel := n.TimesConst(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.PlusEq(rOutParallel)
}

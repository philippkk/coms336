package utils

import (
	"math"
)

// Vec2 represents a 2D vector
type Vec2 struct {
	X, Y float64
}

// ZeroVec2 returns a Vec2 with both components set to 0
func ZeroVec2() Vec2 {
	return Vec2{X: 0, Y: 0}
}

// Length returns the length (magnitude) of the vector
func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// LengthSquared returns the squared length of the vector (no square root)
func (v Vec2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize returns a unit vector pointing in the same direction as v
func (v Vec2) Normalize() Vec2 {
	len := v.Length()
	if len == 0 {
		return ZeroVec2() // return zero vector if length is 0
	}
	return Vec2{X: v.X / len, Y: v.Y / len}
}

// Dot returns the dot product of two Vec2 vectors
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// PlusEq adds another Vec2 to v, modifying v in place
func (v *Vec2) PlusEq(other Vec2) {
	v.X += other.X
	v.Y += other.Y
}

// MinusEq subtracts another Vec2 from v, modifying v in place
func (v *Vec2) MinusEq(other Vec2) {
	v.X -= other.X
	v.Y -= other.Y
}

// TimesConst multiplies v by a constant scalar value, modifying v in place
func (v *Vec2) TimesConst(c float64) {
	v.X *= c
	v.Y *= c
}

// Times returns a new Vec2 that is the result of multiplying v by a scalar
func (v Vec2) Times(c float64) Vec2 {
	return Vec2{X: v.X * c, Y: v.Y * c}
}

// DivideConst divides v by a constant scalar value, modifying v in place
func (v *Vec2) DivideConst(c float64) {
	v.X /= c
	v.Y /= c
}

// Divide returns a new Vec2 that is the result of dividing v by a scalar
func (v Vec2) Divide(c float64) Vec2 {
	return Vec2{X: v.X / c, Y: v.Y / c}
}

// Cross returns the perpendicular vector to v in 2D (only the Z component in 3D is relevant)
func (v Vec2) Cross() float64 {
	// In 2D, cross product is a scalar representing the magnitude in the Z direction.
	return v.X*v.Y - v.Y*v.X
}
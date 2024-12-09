package utils

import (
	"math"
)

type CheckerTexture struct {
	InvScale float64
	Even     Texture
	Odd      Texture
}

// NewCheckerTexture creates a CheckerTexture with two dakar_textures
func NewCheckerTexture(scale float64, even, odd Texture) *CheckerTexture {
	return &CheckerTexture{
		InvScale: 1.0 / scale,
		Even:     even,
		Odd:      odd,
	}
}

// NewCheckerTextureFromColors creates a CheckerTexture with two solid colors
func NewCheckerTextureFromColors(scale float64, c1, c2 Vec3) *CheckerTexture {
	return NewCheckerTexture(scale, NewSolidColor(c1), NewSolidColor(c2))
}

// Value method satisfies the Texture interface for CheckerTexture
func (ct *CheckerTexture) Value(u, v float64, p Vec3) Vec3 {
	xInteger := int(math.Floor(ct.InvScale * p.X))
	yInteger := int(math.Floor(ct.InvScale * p.Y))
	zInteger := int(math.Floor(ct.InvScale * p.Z))

	isEven := (xInteger+yInteger+zInteger)%2 == 0

	if isEven {
		return ct.Even.Value(u, v, p)
	}
	return ct.Odd.Value(u, v, p)
}

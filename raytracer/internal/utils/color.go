package utils

import "math"

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}
func WriteColor(pixels []byte, index int, color Vec3) {
	color = ACESToneMap(color)

	r := LinearToGamma(color.X)
	g := LinearToGamma(color.Y)
	b := LinearToGamma(color.Z)

	intensity := Interval{0.000, 0.999}
	rByte := int(256 * intensity.clamp(r))
	gByte := int(256 * intensity.clamp(g))
	bByte := int(256 * intensity.clamp(b))

	pixels[index] = byte(rByte)
	pixels[index+1] = byte(gByte)
	pixels[index+2] = byte(bByte)
}

func ACESToneMap(color Vec3) Vec3 {
	a := 2.51
	b := 0.03
	c := 2.43
	d := 0.59
	e := 0.14

	color.X = math.Max(0, color.X)
	color.Y = math.Max(0, color.Y)
	color.Z = math.Max(0, color.Z)

	return Vec3{
		X: (color.X * (a*color.X + b)) / (color.X*(c*color.X+d) + e),
		Y: (color.Y * (a*color.Y + b)) / (color.Y*(c*color.Y+d) + e),
		Z: (color.Z * (a*color.Z + b)) / (color.Z*(c*color.Z+d) + e),
	}
}

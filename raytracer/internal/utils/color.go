package utils

import "math"

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}
func WriteColor(pixels []byte, index int, color Vec3) {
	r := color.X
	g := color.Y
	b := color.Z

	r = LinearToGamma(r)
	g = LinearToGamma(g)
	b = LinearToGamma(b)

	intensity := Interval{0.000, 0.999}
	rByte := int(256 * intensity.clamp(r))
	gByte := int(256 * intensity.clamp(g))
	bByte := int(256 * intensity.clamp(b))

	// Directly modify the pre-allocated pixel slice at the specified index
	pixels[index] = byte(rByte)
	pixels[index+1] = byte(gByte)
	pixels[index+2] = byte(bByte)
}
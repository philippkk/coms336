package utils

func WriteColor(pixels *[]byte, color Vec3) {
	r := color.X
	g := color.Y
	b := color.Z

	intensity := Interval{0.000, 0.999}
	rByte := int(256 * intensity.clamp(r))
	gByte := int(256 * intensity.clamp(g))
	bByte := int(256 * intensity.clamp(b))

	*pixels = append(*pixels, byte(rByte), byte(gByte), byte(bByte))
}

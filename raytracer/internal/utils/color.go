package utils

func WriteColor(pixels *[]byte, color Vec3) {
	r := byte(float64(color.X()) * 255.999)
	g := byte(float64(color.Y()) * 255.999)
	b := byte(float64(color.Z()) * 255.999)
	*pixels = append(*pixels, r, g, b)
}

package utils

func WriteColor(pixels *[]byte, color Vec3) {
	r := byte(float64(color.X) * 255.99)
	g := byte(float64(color.Y) * 255.99)
	b := byte(float64(color.Z) * 255.99)
	*pixels = append(*pixels, r, g, b)
}

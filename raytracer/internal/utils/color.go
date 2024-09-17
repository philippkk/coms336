package utils

func WriteColor(pixels *[]byte, color Vec3) {
	r := int((color.X) * 255.99)
	g := int((color.Y) * 255.99)
	b := int((color.Z) * 255.99)
	*pixels = append(*pixels, byte(r), byte(g), byte(b))
}

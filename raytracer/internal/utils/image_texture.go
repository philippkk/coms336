package utils

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strings"
)

type ImageTexture struct {
	Image  *image.RGBA
	Width  int
	Height int
}

// NewImageTexture loads an image file (PNG or JPEG) and creates an ImageTexture
func NewImageTexture(filename string) (*ImageTexture, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Detect file format and decode accordingly
	var img image.Image
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		img, err = png.Decode(file)
	} else if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		img, err = jpeg.Decode(file)
	} else {
		return nil, errors.New("unsupported file format: must be PNG or JPEG")
	}

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// Convert to RGBA with proper color handling
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgba.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	return &ImageTexture{
		Image:  rgba,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}, nil
}

// Clamp clamps a value between a min and max

// Value retrieves the color at the given UV coordinates
func (t *ImageTexture) Value(u, v float64, p Vec3) Vec3 {
	if t.Image == nil || t.Height <= 0 {
		return Vec3{X: 1, Y: 0, Z: 0} // Red to indicate error
	}

	// Print coordinates before wrapping
	originalU, originalV := u, v

	// Handle UV wrapping
	u = math.Mod(u, 1.0)
	v = math.Mod(v, 1.0)
	if u < 0 {
		u += 1.0
	}
	if v < 0 {
		v += 1.0
	}

	// Calculate pixel coordinates
	i := int(u * float64(t.Width-1))
	j := int((1.0 - v) * float64(t.Height-1))

	// Debug output for suspicious values
	if i >= t.Width || j >= t.Height {
		fmt.Printf("Warning: UV mapping out of bounds - Original UV: (%f, %f), "+
			"Wrapped UV: (%f, %f), Pixels: (%d, %d), Dimensions: %dx%d\n",
			originalU, originalV, u, v, i, j, t.Width, t.Height)
	}

	// Ensure we're within bounds
	i = int(math.Max(0, math.Min(float64(i), float64(t.Width-1))))
	j = int(math.Max(0, math.Min(float64(j), float64(t.Height-1))))

	// Get pixel color directly from RGBA buffer
	offset := j*t.Image.Stride + i*4
	r := t.Image.Pix[offset]
	g := t.Image.Pix[offset+1]
	b := t.Image.Pix[offset+2]

	// Debug raw color values for paint texture
	//if t.Width == 768 && t.Height == 768 { // This is the paint texture
	//	fmt.Printf("Paint texture raw colors at UV(%f, %f): RGB(%d, %d, %d)\n",
	//		u, v, r, g, b)
	//}

	// Convert to linear color space with a gentler gamma correction
	colorScale := 1.0 / 255.0
	gamma := 1.8 // Try a lower gamma value

	return Vec3{
		X: math.Pow(float64(r)*colorScale, gamma),
		Y: math.Pow(float64(g)*colorScale, gamma),
		Z: math.Pow(float64(b)*colorScale, gamma),
	}
}

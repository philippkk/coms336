package materials

import (
	"errors"
	"github.com/philippkk/coms336/raytracer/internal/utils"
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
	defer file.Close()

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

	// Convert to RGBA for consistent pixel access
	rgba := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return &ImageTexture{
		Image:  rgba,
		Width:  rgba.Bounds().Dx(),
		Height: rgba.Bounds().Dy(),
	}, nil
}

// Clamp clamps a value between a min and max
func Clamp(x, min, max float64) float64 {
	return math.Max(min, math.Min(x, max))
}

// Value retrieves the color at the given UV coordinates
func (t *ImageTexture) Value(u, v float64, p utils.Vec3) utils.Vec3 {
	// If no texture data, return cyan as a debug color
	if t.Image == nil || t.Height <= 0 {
		return utils.Vec3{X: 0, Y: 1, Z: 1}
	}

	// Clamp UV coordinates to [0,1]
	u = Clamp(u, 0.0, 1.0)
	v = 1.0 - Clamp(v, 0.0, 1.0) // Flip V to image coordinates

	// Map UV to pixel coordinates
	i := int(u * float64(t.Width))
	j := int(v * float64(t.Height))

	// Clamp indices to valid range
	if i >= t.Width {
		i = t.Width - 1
	}
	if j >= t.Height {
		j = t.Height - 1
	}

	// Get pixel color
	pixel := t.Image.At(i, j).(color.RGBA)

	// Scale color to [0,1]
	colorScale := 1.0 / 255.0
	return utils.Vec3{
		X: float64(pixel.R) * colorScale,
		Y: float64(pixel.G) * colorScale,
		Z: float64(pixel.B) * colorScale,
	}
}

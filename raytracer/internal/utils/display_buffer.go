package utils

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"image"
	"image/color"
)

// DisplayBuffer represents the window and its associated state
type DisplayBuffer struct {
	Win    *opengl.Window
	canvas *image.RGBA
	pic    *pixel.PictureData
}

// NewDisplayBuffer creates a new window for displaying the raytracer output
func NewDisplayBuffer(width, height int) (*DisplayBuffer, error) {
	cfg := opengl.WindowConfig{
		Title:  "cool raytracer",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	pic := pixel.PictureDataFromImage(canvas)
	win.SetSmooth(true)
	return &DisplayBuffer{
		Win:    win,
		canvas: canvas,
		pic:    pic,
	}, nil
}

func (d *DisplayBuffer) UpdatePixel(x, y int, col color.Color) {
	d.canvas.Set(x, y, col)
}

func (d *DisplayBuffer) Refresh() {
	d.pic = pixel.PictureDataFromImage(d.canvas)
	sprite := pixel.NewSprite(d.pic, d.pic.Bounds())

	//d.Win.Clear(colornames.Black)
	sprite.Draw(d.Win, pixel.IM.Moved(d.Win.Bounds().Center()))
	d.Win.Update()
}

func (d *DisplayBuffer) ShouldClose() bool {
	return d.Win.Closed()
}

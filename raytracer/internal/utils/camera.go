package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
)

type Camera struct {
	AspectRatio                                  float64
	ImageWidth, imageHeight                      int
	center, pixel00Loc, pixelDeltaU, pixelDeltaV Vec3
}

func (c *Camera) Render(world HittableList) {
	c.initialize()
	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", c.ImageWidth, c.imageHeight, 255)

	var pixels []byte
	for j := 0; j < c.imageHeight; j++ {
		for i := 0; i < c.ImageWidth; i++ {
			pixelCenter := c.pixel00Loc.PlusEq(c.pixelDeltaU.TimesConst(float64(i))).PlusEq(c.pixelDeltaV.TimesConst(float64(j)))
			rayDirection := pixelCenter.MinusEq(c.center)

			ray := Ray{Origin: c.center, Direction: rayDirection}

			color := rayColor(ray, &world)
			WriteColor(&pixels, color)
		}
	}

	_, err = file.Write(pixels)
	if err != nil {
		return
	}

	fmt.Println("Done.")

	openFile("goimage.ppm")
}
func (c *Camera) initialize() {
	c.imageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.imageHeight < 0 {
		c.imageHeight = 1
	}
	c.center = Vec3{0, 0, 0}

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(c.ImageWidth) / float64(c.imageHeight))
	cameraCenter := Vec3{}

	viewportU := Vec3{X: viewportWidth}
	viewportV := Vec3{Y: -viewportHeight}
	c.pixelDeltaU = viewportU.TimesConst(1.0 / float64(c.ImageWidth))
	c.pixelDeltaV = viewportV.TimesConst(1.0 / float64(c.imageHeight))

	viewportUpperLeft := cameraCenter.MinusEq(Vec3{Z: focalLength}).MinusEq(viewportU.TimesConst(0.5)).MinusEq(viewportV.TimesConst(0.5))
	c.pixel00Loc = c.pixelDeltaU.PlusEq(c.pixelDeltaV).TimesConst(0.5).PlusEq(viewportUpperLeft)
}
func rayColor(r Ray, world Hittable) Vec3 {
	var rec HitRecord
	if world.Hit(&r, Interval{0, math.Inf(+1)}, &rec) {
		temp := rec.Normal.PlusEq(Vec3{X: 1, Y: 1, Z: 1})
		return temp.TimesConst(0.5)
	}

	unitDirection := r.Direction.Normalize()
	a := 0.5 * (unitDirection.Y + 1.0)
	white := Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue := Vec3{X: 0.5, Y: 0.7, Z: 1.0}
	return white.TimesConst(1.0 - a).PlusEq(blue.TimesConst(a))
}

func openFile(filename string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "darwin":
		cmd = exec.Command("open", filename)
	default:
		cmd = exec.Command("xdg-open", filename)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
}
package main

import (
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
	"os"
	"os/exec"
	"runtime"
)

func HitSphere(r utils.Ray, s utils.Vec3, radius float64) float64 {
	oc := s.MinusEq(r.Origin)
	a := r.Direction.LengthSquared()
	h := r.Direction.Dot(oc)
	c := oc.LengthSquared() - radius*radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return -1.0
	} else {
		return (h - math.Sqrt(discriminant)) / a
	}
}
func rayColor(r utils.Ray) utils.Vec3 {
	t := HitSphere(r, utils.Vec3{Z: -1}, 0.5)
	if t > 0 {
		N := r.At(t).MinusEq(utils.Vec3{Z: -1}).UnitVector()
		return utils.Vec3{X: N.X + 1, Y: N.Y + 1, Z: N.Z + 1}.TimesConst(0.5)
	}

	unitDirection := r.Direction.Normalize()
	a := 0.5 * (unitDirection.Y + 1.0)
	white := utils.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue := utils.Vec3{X: 0.5, Y: 0.7, Z: 1.0}
	return white.TimesConst(1.0 - a).PlusEq(blue.TimesConst(a))
}
func main() {
	aspectRatio := 16.0 / 9.0
	width := 800
	height := int(float64(width) / aspectRatio)
	if height < 0 {
		height = 1
	}

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(width) / float64(height))
	cameraCenter := utils.Vec3{}

	viewportU := utils.Vec3{X: viewportWidth}
	viewportV := utils.Vec3{Y: -viewportHeight}
	pixelDeltaU := viewportU.TimesConst(1.0 / float64(width))
	pixelDeltaV := viewportV.TimesConst(1.0 / float64(height))

	viewportUpperLeft := cameraCenter.MinusEq(utils.Vec3{Z: focalLength}).MinusEq(viewportU.TimesConst(0.5)).MinusEq(viewportV.TimesConst(0.5))
	pixel00Loc := pixelDeltaU.PlusEq(pixelDeltaV).TimesConst(0.5).PlusEq(viewportUpperLeft)

	maxColorValue := 255

	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", width, height, maxColorValue)

	var pixels []byte
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			pixelCenter := pixel00Loc.PlusEq(pixelDeltaU.TimesConst(float64(i))).PlusEq(pixelDeltaV.TimesConst(float64(j)))
			rayDirection := pixelCenter.MinusEq(cameraCenter)

			ray := utils.Ray{Origin: cameraCenter, Direction: rayDirection}

			color := rayColor(ray)
			utils.WriteColor(&pixels, color)
		}
	}

	_, err = file.Write(pixels)
	if err != nil {
		return
	}

	fmt.Println("Done.")

	test := utils.Sphere{5}
	test2 := utils.Tri{6}
	utils.Measure(test)
	utils.Measure(test2)

	openFile("goimage.ppm")
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

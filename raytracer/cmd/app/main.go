package main

import (
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
	"os"
	"os/exec"
	"runtime"
)

func rayColor(r utils.Ray, world utils.Hittable) utils.Vec3 {
	var rec utils.HitRecord
	if world.Hit(&r, 0, math.Inf(1), &rec) {
		fmt.Println(rec.Normal)
		temp := rec.Normal.PlusEq(utils.Vec3{X: 1, Y: 1, Z: 1})
		return temp.TimesConst(0.5)
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

	var world utils.HittableList
	world.Add(utils.Sphere{utils.Vec3{0, 0, -1}, 0.5})
	world.Add(utils.Sphere{utils.Vec3{0, -100.5, -1}, 100})

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

			color := rayColor(ray, &world)
			utils.WriteColor(&pixels, color)
		}
	}

	_, err = file.Write(pixels)
	if err != nil {
		return
	}

	fmt.Println("Done.")

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

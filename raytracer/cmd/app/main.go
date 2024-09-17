package main

import (
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"os"
	"os/exec"
	"runtime"
)

func rayColor(r utils.Ray) utils.Vec3 {
	unitDirection := r.Direction.UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)

	white := utils.NewVec3(1.0, 1.0, 1.0)
	blue := utils.NewVec3(0.5, 0.7, 1.0)
	timesConst := white.TimesConst(1.0 - a)
	return timesConst.PlusEq(blue.TimesConst(a))
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
	cameraCenter := utils.NewVec3(0, 0, 0)

	viewportU := utils.NewVec3(viewportWidth, 0, 0)
	viewportV := utils.NewVec3(0, -viewportHeight, 0)

	pixelDeltaU := viewportU.TimesConst(float64(1 / width))
	pixelDeltaV := viewportV.TimesConst(float64(1 / height))

	viewportUpperLeft := cameraCenter.MinusEq(*utils.NewVec3(0, 0, focalLength))
	viewportUpperLeft = viewportUpperLeft.MinusEq(viewportU.TimesConst(1.0 / 2.0))
	viewportUpperLeft = viewportUpperLeft.MinusEq(viewportV.TimesConst(1.0 / 2.0))

	plusConst := viewportUpperLeft.PlusConst(0.5)
	pixel00Loc := plusConst.TimesEq(pixelDeltaU.PlusEq(pixelDeltaV))

	maxColorValue := 255

	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", width, height, maxColorValue)

	var pixels []byte
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			eq := pixel00Loc.PlusEq(pixelDeltaU.TimesConst(float64(j)))
			pixelCenter := eq.PlusEq(pixelDeltaV.TimesConst(float64(i)))
			rayDirection := pixelCenter.MinusEq(*cameraCenter)

			fmt.Println(rayDirection)
			ray := utils.Ray{Origin: *cameraCenter, Direction: rayDirection}

			color := rayColor(ray)
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

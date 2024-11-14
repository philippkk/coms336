package main

import (
	"fmt"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/philippkk/coms336/raytracer/internal/materials"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var t time.Duration

func run() {
	// Your existing main() code goes here
	// world setup, camera setup, etc.
	var world utils.HittableList

	//materialGround := materials.Lambertian{utils.Vec3{0.07, 0.2, 0.05}}
	materialLeft := materials.Lambertian{utils.Vec3{0.0, 0.3, 0.7}}
	materialLeft2 := materials.Lambertian{utils.Vec3{0.7, 0.0, 0.0}}
	materialRight := materials.Metal{utils.Vec3{0.0, 0.3, 0.7}, 0.0}
	materialGlass := materials.Dielectric{1.50}
	materialGlass2 := materials.Dielectric{1.00 / 1.50}

	world.Add(objects.Triangle{
		utils.Vec3{1, 0, -1.0},
		utils.Vec3{0, 0, -1.0},
		utils.Vec3{0, 0.5, 0},
		materialLeft2})
	//world.Add(objects.Sphere{utils.Ray{utils.Vec3{2, 0, -1}, utils.Vec3{0, 0, 0}, 0}, 1, materialGlass})
	world.Add(objects.Sphere{utils.Ray{utils.Vec3{2, 0, -1}, utils.Vec3{0, -0.2, 0}, 0}, 0.5, materialRight})
	world.Add(objects.Sphere{utils.Ray{utils.Vec3{-1, 0, -1}, utils.Vec3{0, 0, 0}, 0}, 0.5, materialGlass})
	world.Add(objects.Sphere{utils.Ray{utils.Vec3{-1, 0, -1}, utils.Vec3{0, 0, 0}, 0}, 0.4, materialGlass2})
	world.Add(objects.Sphere{utils.Ray{utils.Vec3{-1, 0, -1}, utils.Vec3{0, 0, 0}, 0}, 0.2, materialLeft})
	world.Add(objects.Sphere{utils.Ray{utils.Vec3{0, -100.5, -1}, utils.Vec3{0, 0, 0}, 0}, 100, materialRight})

	var cam utils.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 800 //2234
	cam.SamplesPerPixel = 100
	cam.MaxDepth = 100

	cam.Vfov = 30
	cam.LookFrom = utils.Vec3{0, 1, -6}
	cam.LookAt = utils.Vec3{0, 0, 0}
	cam.Vup = utils.Vec3{Y: 1}
	cam.DefocusAngle = 0.0 //0.6 was nice;
	cam.Focusdist = 5

	imageHeight := int(float64(cam.ImageWidth) / cam.AspectRatio)
	if imageHeight < 0 {
		imageHeight = 1
	}
	// Create display buffer
	display, err := utils.NewDisplayBuffer(cam.ImageWidth, imageHeight)
	if err != nil {
		fmt.Println("Error creating display:", err)
		return
	}

	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", cam.ImageWidth, imageHeight, 255)
	pixels := make([]byte, imageHeight*cam.ImageWidth*3)

	for !display.Win.Closed() {
		if display.Win.Pressed(pixel.KeyW) {
			cam.LookFrom.Z += 1
		}
		if display.Win.Pressed(pixel.KeyS) {
			cam.LookFrom.Z -= 1
		}
		if display.Win.Pressed(pixel.KeyD) {
			cam.LookFrom.X -= 1
		}
		if display.Win.Pressed(pixel.KeyA) {
			cam.LookFrom.X += 1
		}
		finishTime := cam.Render(world, display, pixels)
		if finishTime > t {
			t = finishTime
		}
		display.Win.Update()
	}

	/*
		todo: seems to write after the first worker finished instead of the last
	*/
	if display.ShouldClose() {
		_, err = file.Write(pixels)
		if err != nil {
			return
		}

		fmt.Printf("\033[1A\033[K")
		fmt.Printf("Max image time: %v\n", t)
		fmt.Printf("Image size: %d x %d\n", cam.ImageWidth, imageHeight)

		//openFile("goimage.ppm")
	}

}
func main() {
	opengl.Run(run)
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

package main

import (
	"fmt"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/philippkk/coms336/raytracer/internal/materials"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
)

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
	world.Add(objects.Sphere{utils.Vec3{2, 0, -1}, 1, materialGlass})
	world.Add(objects.Sphere{utils.Vec3{2, 0, -1}, 0.5, materialRight})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.5, materialGlass})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.4, materialGlass2})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.2, materialLeft})
	//
	world.Add(objects.Sphere{utils.Vec3{0, -100.5, -1}, 100, materialRight})
	//
	//

	var cam utils.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1200 //2234
	cam.SamplesPerPixel = 200
	cam.MaxDepth = 200

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

	cam.Render(world, display)
	for !display.Win.Closed() {
		display.Win.Update()
	}

}
func main() {
	opengl.Run(run)
}

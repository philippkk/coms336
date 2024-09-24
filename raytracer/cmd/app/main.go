package main

import (
	"github.com/philippkk/coms336/raytracer/internal/materials"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
)

func main() {

	var world utils.HittableList

	materialGround := materials.Lambertian{utils.Vec3{0.4, 0.0, 0.7}}
	//materialLeft := materials.Metal{utils.Vec3{0.0, 0.3, 0.7}, 1.0}
	//materialRight := materials.Metal{utils.Vec3{0.8, 0.8, 0.8}, 0.0}
	materialGlass := materials.Dielectric{1.00 / 1.50}

	world.Add(objects.Sphere{utils.Vec3{0, 0, -1}, 0.5, materialGlass})
	world.Add(objects.Sphere{utils.Vec3{1, 0, -1}, 0.2, materialGround})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.2, materialGround})

	world.Add(objects.Sphere{utils.Vec3{0, -100.5, -1}, 100, materialGround})

	var cam utils.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1920
	cam.SamplesPerPixel = 10
	cam.MaxDepth = 10
	cam.Vfov = 90

	cam.Render(world)
}

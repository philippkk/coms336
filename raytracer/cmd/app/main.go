package main

import (
	"github.com/philippkk/coms336/raytracer/internal/materials"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
)

func main() {

	var world utils.HittableList

	//materialGround := materials.Lambertian{utils.Vec3{0.07, 0.2, 0.05}}
	materialLeft := materials.Lambertian{utils.Vec3{0.0, 0.3, 0.7}}
	materialRight := materials.Metal{utils.Vec3{0.0, 0.3, 0.7}, 0.0}
	materialGlass := materials.Dielectric{1.50}
	materialGlass2 := materials.Dielectric{1.00 / 1.50}
	world.Add(objects.Sphere{utils.Vec3{0, 0, -1.2}, 0.5, materialRight})
	world.Add(objects.Sphere{utils.Vec3{2, 0, -1}, 1, materialGlass})
	world.Add(objects.Sphere{utils.Vec3{2, 0, -1}, 0.5, materialLeft})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.5, materialGlass})
	world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.4, materialGlass2})
	//world.Add(objects.Sphere{utils.Vec3{-1, 0, -1}, 0.2, materialLeft})

	world.Add(objects.Sphere{utils.Vec3{0, -100.5, -1}, 100, materialLeft})

	var cam utils.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 2234
	cam.SamplesPerPixel = 200
	cam.MaxDepth = 10000
	cam.Vfov = 30
	cam.LookFrom = utils.Vec3{-2, 1.5, 1}
	cam.LookAt = utils.Vec3{Z: -1}
	cam.Vup = utils.Vec3{Y: 1}
	cam.DefocusAngle = 0.6
	cam.Focusdist = 3.4

	cam.Render(world)
}

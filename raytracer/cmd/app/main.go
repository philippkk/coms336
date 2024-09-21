package main

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
)

func main() {

	var world utils.HittableList
	world.Add(utils.Sphere{utils.Vec3{0, 0, -1}, 0.5})
	world.Add(utils.Sphere{utils.Vec3{1, 0, -1}, 0.2})
	world.Add(utils.Sphere{utils.Vec3{-1, 0, -1}, 0.2})

	world.Add(utils.Sphere{utils.Vec3{0, -100.5, -1}, 100})

	var cam utils.Camera
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1920

	cam.Render(world)
}

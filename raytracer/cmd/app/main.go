package main

import (
	"fmt"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/philippkk/coms336/raytracer/internal/materials"
	model2 "github.com/philippkk/coms336/raytracer/internal/model"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var t time.Duration
var world utils.HittableList
var cam utils.Camera

var reRender bool

func run() {
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

	moveAmount := 10.0
	for !display.Win.Closed() {
		if display.Win.Pressed(pixel.KeyW) {
			cam.LookFrom.Z += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyS) {
			cam.LookFrom.Z -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyD) {
			cam.LookFrom.X -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyA) {
			cam.LookFrom.X += moveAmount
			reRender = true
		}
		if reRender {
			finishTime := cam.Render(world, display, pixels)
			if finishTime > t {
				t = finishTime
			}
			reRender = false
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
	world = createRandomScene()

	fmt.Printf("\033[1A\033[K")
	fmt.Println("loadin file")
	model := model2.NewModel("internal/model/smolengine.obj")
	fmt.Printf("\033[1A\033[K")
	fmt.Printf("done!")

	triangleMat := materials.Metal{utils.Vec3{1, 0, 0}, 0.2}

	triangles := model.ToTriangles(triangleMat)
	for _, triangle := range triangles {
		world.Add(triangle)
	}

	fmt.Println("\n num of objects: ", len(world.Objects))
	fmt.Println(" ")
	bvhRoot := utils.NewBVHNode(world.Objects, 0, len(world.Objects))
	world = utils.HittableList{Objects: []utils.Hittable{bvhRoot}}

	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1000 //2234
	cam.SamplesPerPixel = 10
	cam.MaxDepth = 10

	cam.Vfov = 30
	cam.LookFrom = utils.Vec3{3, 3, 7}
	cam.LookAt = utils.Vec3{0, 0, 0}
	cam.Vup = utils.Vec3{Y: 1}
	cam.DefocusAngle = 0.0 //0.6 was nice;
	cam.Focusdist = 5

	reRender = true
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

func createRandomScene() utils.HittableList {
	var world utils.HittableList

	groundMaterial := materials.Lambertian{Albedo: utils.Vec3{0.5, 0.5, 0.5}}
	world.Add(objects.CreateSphere(
		utils.Ray{Origin: utils.Vec3{0, -1000, 0}},
		1000,
		groundMaterial))

	num := 10
	// Add many small random spheres
	for a := -num; a < num; a++ {
		for b := -num; b < num; b++ {
			chooseMat := rand.Float64()
			center := utils.Vec3{
				X: float64(a) + 0.9*rand.Float64(),
				Y: 0.2,
				Z: float64(b) + 0.9*rand.Float64(),
			}

			if center.MinusEq(utils.Vec3{4, 0.2, 0}).Length() > 0.9 {
				var sphereMaterial utils.Material

				if chooseMat < 0.8 {
					// Diffuse
					albedo := utils.Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
					sphereMaterial = materials.Lambertian{Albedo: albedo}
				} else if chooseMat < 0.95 {
					// Metal
					albedo := utils.Vec3{
						rand.Float64()*0.5 + 0.5,
						rand.Float64()*0.5 + 0.5,
						rand.Float64()*0.5 + 0.5,
					}
					fuzz := rand.Float64() * 0.5
					sphereMaterial = materials.Metal{Albedo: albedo, Fuzz: fuzz}
				} else {
					// Glass
					sphereMaterial = materials.Dielectric{RefractionIndex: 1.5}
				}

				world.Add(objects.CreateSphere(
					utils.Ray{Origin: center},
					0.2,
					sphereMaterial))
			}
		}
	}

	return world
}

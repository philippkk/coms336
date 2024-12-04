package main

import (
	"fmt"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/philippkk/coms336/raytracer/internal/materials"
	"github.com/philippkk/coms336/raytracer/internal/model"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"image"
	"image/color"
	"image/png"
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

	//file, err := os.Create("goimage.ppm")
	//if err != nil {
	//	fmt.Println("Error creating file:", err)
	//	return
	//}
	//defer file.Close()

	//fmt.Fprintf(file, "P6\n%d %d\n%d\n", cam.ImageWidth, imageHeight, 255)
	pixels := make([]byte, imageHeight*cam.ImageWidth*3)

	moveAmount := 0.5
	for !display.Win.Closed() {
		if display.Win.Pressed(pixel.KeyQ) {
			cam.LookFrom.X += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyA) {
			cam.LookFrom.X -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyW) {
			cam.LookFrom.Y += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyS) {
			cam.LookFrom.Y -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyE) {
			cam.LookFrom.Z += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyD) {
			cam.LookFrom.Z -= moveAmount
			reRender = true
		}

		if display.Win.Pressed(pixel.KeyU) {
			cam.LookAt.X += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyJ) {
			cam.LookAt.X -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyI) {
			cam.LookAt.Y += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyK) {
			cam.LookAt.Y -= moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyO) {
			cam.LookAt.Z += moveAmount
			reRender = true
		}
		if display.Win.Pressed(pixel.KeyL) {
			cam.LookAt.Z -= moveAmount
			reRender = true
		}

		if reRender {
			finishTime := cam.Render(world, display, pixels)
			if finishTime > t {
				t = finishTime
			}
			reRender = false
			println("from:", cam.LookFrom.X, cam.LookFrom.Y, cam.LookFrom.Z)
			println("at:", cam.LookAt.X, cam.LookAt.Y, cam.LookAt.Z)
			println(" ")
		}
		display.Win.Update()
	}

	/*
		todo: seems to write after the first worker finished instead of the last
	*/
	fmt.Printf("\033[1A\033[K")
	fmt.Printf("Max image time: %v\n", t)
	fmt.Printf("Image size: %d x %d\n", cam.ImageWidth, imageHeight)

	saveToPNG(cam.ImageWidth, imageHeight, pixels)
	return
	if display.ShouldClose() {
		//_, err = file.Write(pixels)
		//if err != nil {
		//	return
		//}

		fmt.Printf("\033[1A\033[K")
		fmt.Printf("Max image time: %v\n", t)
		fmt.Printf("Image size: %d x %d\n", cam.ImageWidth, imageHeight)

		saveToPNG(cam.ImageWidth, imageHeight, pixels)
		//openFile("goimage.ppm")
	}

}
func main() {
	world = createRandomScene()
	//
	fmt.Printf("\033[1A\033[K")
	fmt.Println("loadin file")
	newmod := model.NewModel("internal/model/dakar.obj", "internal/model/dakar.mtl")
	fmt.Printf("\033[1A\033[K")
	fmt.Printf("done!")
	////
	earthmap, _ := materials.NewImageTexture("internal/model/earthmap.jpg")
	triangleMat := materials.NewLambertian(earthmap)
	//triangleMat := materials.Metal{utils.Vec3{1, 0, 0}, 0.0}
	//
	triangles := newmod.ToTriangles(triangleMat, false)
	for _, triangle := range triangles {
		world.Add(triangle)
	}

	//testTri := objects.CreateTriangle(
	//	utils.Vec3{1, 2, 0},
	//	utils.Vec3{2, 2, 0},
	//	utils.Vec3{2, 3, 0},
	//	triangleMat)
	//world.Add(testTri)

	fmt.Println("\n num of objects: ", len(world.Objects))
	fmt.Println(" ")
	bvhRoot := utils.NewBVHNode(world.Objects, 0, len(world.Objects))
	world = utils.HittableList{Objects: []utils.Hittable{bvhRoot}}

	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 800 //2234
	cam.SamplesPerPixel = 2
	cam.MaxDepth = 2

	cam.Vfov = 30
	cam.LookFrom = utils.Vec3{5, 5, 5}
	cam.LookAt = utils.Vec3{-2, 0, -1}
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

	checker := materials.NewCheckerTextureFromColors(0.64,
		utils.Vec3{0.0940, 0.078, 0.149},
		utils.Vec3{0.596, 0.435, 0.839})
	//earthmap, _ := materials.NewImageTexture("internal/model/earthmap.jpg")
	groundMaterial := materials.NewLambertian(checker)
	world.Add(objects.CreateSphere(
		utils.Ray{Origin: utils.Vec3{0, -1000, 0}},
		1000,
		groundMaterial))

	num := 0
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
					sphereMaterial = materials.NewLambertianFromColor(albedo)
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

func saveToPNG(width, height int, pixels []byte) {
	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the RGBA image with your pixel data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 3
			r := pixels[index]
			g := pixels[index+1]
			b := pixels[index+2]

			// Set pixel color in the RGBA image
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// Save to PNG file
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode as PNG
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	println("PNG file saved as output.png")
}

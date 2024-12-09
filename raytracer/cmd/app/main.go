package main

import (
	"fmt"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/philippkk/coms336/raytracer/internal/model"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"github.com/philippkk/coms336/raytracer/internal/utils/material"
	"image"
	"image/color"
	"image/png"
	"math"
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

type Scene struct {
	World          utils.HittableList
	Cam            utils.Camera
	CubeMap        *utils.CubeMap
	SkipBackground bool
}

func createDefaultCam() utils.Camera {
	return utils.Camera{
		AspectRatio:     16.0 / 9.0,
		ImageWidth:      1200,
		SamplesPerPixel: 200,
		MaxDepth:        50,
		Vfov:            30,
		DefocusAngle:    0.0,
		Focusdist:       5,
	}
}
func CreateCornellBox() Scene {
	var config Scene

	// Materials
	red := material.NewLambertianFromColor(utils.Vec3{X: 0.65, Y: 0.05, Z: 0.05})
	white := material.NewLambertianFromColor(utils.Vec3{X: 0.73, Y: 0.73, Z: 0.73})
	green := material.NewLambertianFromColor(utils.Vec3{X: 0.12, Y: 0.45, Z: 0.15})
	light := material.NewDiffuseLightFromColor(utils.Vec3{X: 15, Y: 15, Z: 15})

	// Room walls
	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 555, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 555, Z: 0},
		utils.Vec3{X: 0, Y: 0, Z: 555},
		green,
	))

	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 0, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 555, Z: 0},
		utils.Vec3{X: 0, Y: 0, Z: 555},
		red,
	))

	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 343, Y: 554, Z: 332},
		utils.Vec3{X: -130, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 0, Z: -105},
		light,
	))

	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 0, Y: 0, Z: 0},
		utils.Vec3{X: 555, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 0, Z: 555},
		white,
	))

	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 555, Y: 555, Z: 555},
		utils.Vec3{X: -555, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 0, Z: -555},
		white,
	))

	config.World.Add(objects.CreateQuad(
		utils.Vec3{X: 0, Y: 0, Z: 555},
		utils.Vec3{X: 555, Y: 0, Z: 0},
		utils.Vec3{X: 0, Y: 555, Z: 0},
		white,
	))

	for _, box := range objects.CreateBox(
		utils.Vec3{130, 0, 65},
		utils.Vec3{295, 165, 230},
		white) {
		config.World.Add(box)
	}

	for _, box := range objects.CreateBox(
		utils.Vec3{265, 0, 295},
		utils.Vec3{430, 330, 460},
		white) {
		config.World.Add(box)
	}

	// Camera settings
	config.Cam = utils.Camera{
		AspectRatio:     1.0,
		ImageWidth:      600,
		SamplesPerPixel: 200,
		MaxDepth:        50,
		Vfov:            40,
		LookFrom:        utils.Vec3{X: 278, Y: 278, Z: -800},
		LookAt:          utils.Vec3{X: 278, Y: 278, Z: 0},
		Vup:             utils.Vec3{Y: 1},
		DefocusAngle:    0,
		Focusdist:       10.0,
	}

	// Scene settings
	config.SkipBackground = true

	return config
}

func createQuadsScene() Scene {
	var scene Scene

	// Materials
	leftRed := material.NewLambertianFromColor(utils.Vec3{X: 1.0, Y: 0.2, Z: 0.2})
	backGreen := material.NewLambertianFromColor(utils.Vec3{X: 0.2, Y: 1.0, Z: 0.2})
	rightBlue := material.NewLambertianFromColor(utils.Vec3{X: 0.2, Y: 0.2, Z: 1.0})
	upperOrange := material.NewLambertianFromColor(utils.Vec3{X: 1.0, Y: 0.5, Z: 0.0})
	lowerTeal := material.NewLambertianFromColor(utils.Vec3{X: 0.2, Y: 0.8, Z: 0.8})

	// Create quads
	scene.World.Add(objects.CreateQuad(
		utils.Vec3{X: -3, Y: -2, Z: 5}, // point
		utils.Vec3{X: 0, Y: 0, Z: -4},  // u - towards back
		utils.Vec3{X: 0, Y: 4, Z: 0},   // v - up
		leftRed,
	))

	scene.World.Add(objects.CreateQuad(
		utils.Vec3{X: -2, Y: -2, Z: 0}, // point
		utils.Vec3{X: 4, Y: 0, Z: 0},   // u - right
		utils.Vec3{X: 0, Y: 4, Z: 0},   // v - up
		backGreen,
	))

	scene.World.Add(objects.CreateQuad(
		utils.Vec3{X: 3, Y: -2, Z: 1}, // point
		utils.Vec3{X: 0, Y: 0, Z: 4},  // u - forward
		utils.Vec3{X: 0, Y: 4, Z: 0},  // v - up
		rightBlue,
	))

	scene.World.Add(objects.CreateQuad(
		utils.Vec3{X: -2, Y: 3, Z: 1}, // point
		utils.Vec3{X: 4, Y: 0, Z: 0},  // u - right
		utils.Vec3{X: 0, Y: 0, Z: 4},  // v - forward
		upperOrange,
	))

	scene.World.Add(objects.CreateQuad(
		utils.Vec3{X: -2, Y: -3, Z: 5}, // point
		utils.Vec3{X: 4, Y: 0, Z: 0},   // u - right
		utils.Vec3{X: 0, Y: 0, Z: -4},  // v - back
		lowerTeal,
	))

	// Camera settings
	scene.Cam = utils.Camera{
		AspectRatio:     1.0,
		ImageWidth:      400,
		SamplesPerPixel: 100,
		MaxDepth:        50,
		Vfov:            80,
		LookFrom:        utils.Vec3{X: 0, Y: 0, Z: 9},
		LookAt:          utils.Vec3{X: 0, Y: 0, Z: 0},
		Vup:             utils.Vec3{Y: 1},
		DefocusAngle:    0,
		Focusdist:       10.0,
	}

	// Set black background
	scene.SkipBackground = false

	return scene
}
func createRandomScene() Scene {
	var scene Scene
	scene.World = utils.HittableList{}

	num := 2
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
					sphereMaterial = material.NewLambertianFromColor(albedo)
				} else if chooseMat < 0.95 {
					// Metal
					albedo := utils.Vec3{
						rand.Float64()*0.5 + 0.5,
						rand.Float64()*0.5 + 0.5,
						rand.Float64()*0.5 + 0.5,
					}
					fuzz := rand.Float64() * 0.5
					sphereMaterial = material.Metal{Albedo: albedo, Fuzz: fuzz}
				} else {
					// Glass
					sphereMaterial = material.Dielectric{RefractionIndex: 1.5}
				}

				//scene.World.Add(objects.CreateSphere(
				//	utils.Ray{Origin: center},
				//	0.2,
				//	sphereMaterial))

				// SMOKE BALLS!
				scene.World.Add(material.CreateConstantMedium(
					objects.CreateSphere(
						utils.Ray{Origin: center},
						0.2,
						sphereMaterial),
					.8,
					utils.Vec3{12, 0, 0}))
			}
		}
	}

	// Random scene specific camera settings
	scene.Cam = createDefaultCam()
	scene.Cam.LookFrom = utils.Vec3{4.5, 2.5, 3.5}
	scene.Cam.LookAt = utils.Vec3{-2, 0, -1}
	scene.Cam.Vup = utils.Vec3{Y: 1}

	return scene
}

func createModelScene() Scene {
	var scene Scene

	fmt.Println("Loading file")
	newmod := model.NewModel("internal/model/dakar.obj", "internal/model/dakar.mtl")
	fmt.Printf("done!")

	defaultMat := material.Dielectric{RefractionIndex: 1.520}

	triangles := newmod.ToTriangles(defaultMat, "dakar_textures")
	for _, triangle := range triangles {
		scene.World.Add(triangle)
	}

	//newmod2 := model.NewModel("internal/model/mine.obj", "internal/model/mine.mtl")
	//
	//triangles = newmod2.ToTriangles(defaultMat, "minecraft_textures")
	//for _, triangle := range triangles {
	//	world.Add(triangle)
	//}

	//newmod3 := model.NewModel("internal/model/rauh.obj", "internal/model/rauh.mtl")
	//
	//triangles = newmod3.ToTriangles(defaultMat, "rauh_textures")
	//for _, triangle := range triangles {
	//	world.Add(triangle)
	//}

	mat := material.NewDiffuseLightFromColor(utils.Vec3{12, 0, 0})
	scene.World.Add(objects.CreateSphere(
		utils.Ray{Origin: utils.Vec3{0, 0, 0}},
		1,
		mat))

	// Model scene specific camera settings
	scene.Cam = createDefaultCam()
	scene.Cam.LookFrom = utils.Vec3{4.5, 2.5, 3.5}
	scene.Cam.LookAt = utils.Vec3{-2, 0, -1}
	scene.Cam.Vup = utils.Vec3{Y: 1}

	return scene
}

func createQuadricScene() Scene {
	var scene Scene

	scene.World.Add(objects.CreateQuadricSphere(
		utils.Vec3{X: 0, Y: 0, Z: 0},
		1.0,
		material.NewLambertianFromColor(utils.Vec3{X: 0.7, Y: 0.3, Z: 0.3}),
	))

	scene.World.Add(objects.CreateCylinder(
		utils.Vec3{X: 0, Y: 0, Z: 0},
		0.5,
		material.NewLambertianFromColor(utils.Vec3{X: 0.3, Y: 0.7, Z: 0.3}),
	))

	scene.World.Add(objects.CreateCone(
		utils.Vec3{X: -2, Y: 0, Z: 0},
		math.Pi/6, // 30 degrees
		material.NewLambertianFromColor(utils.Vec3{X: 0.3, Y: 0.3, Z: 0.7}),
	))

	mat := material.NewDiffuseLightFromColor(utils.Vec3{12, 0, 0})
	scene.World.Add(objects.CreateSphere(
		utils.Ray{Origin: utils.Vec3{0, 0, 0}},
		1,
		mat))

	// Model scene specific camera settings
	scene.Cam = createDefaultCam()
	scene.Cam.LookFrom = utils.Vec3{4.5, 2.5, 3.5}
	scene.Cam.LookAt = utils.Vec3{-2, 0, -1}
	scene.Cam.Vup = utils.Vec3{Y: 1}

	return scene
}

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
	//fmt.Printf("\033[1A\033[K")
	fmt.Printf("Max image time: %v\n", t)
	fmt.Printf("Image size: %d x %d\n", cam.ImageWidth, imageHeight)

	saveToPNG(cam.ImageWidth, imageHeight, pixels)
	return
	if display.ShouldClose() {
		//_, err = file.Write(pixels)
		//if err != nil {
		//	return
		//}

		//fmt.Printf("\033[1A\033[K")
		fmt.Printf("Max image time: %v\n", t)
		fmt.Printf("Image size: %d x %d\n", cam.ImageWidth, imageHeight)

		saveToPNG(cam.ImageWidth, imageHeight, pixels)
		//openFile("goimage.ppm")
	}

}
func main() {
	scene := CreateCornellBox()

	fmt.Println("\n num of objects: ", len(scene.World.Objects))
	fmt.Println(" ")

	cubeMap, err := utils.NewCubeMap(
		"internal/utils/cube_map_images/posx.jpg", // RIGHT
		"internal/utils/cube_map_images/negx.jpg", // LEFT
		"internal/utils/cube_map_images/posy.jpg", // TOP
		"internal/utils/cube_map_images/negy.jpg", // BOTTOM
		"internal/utils/cube_map_images/posz.jpg", // FRONT
		"internal/utils/cube_map_images/negz.jpg", // BACK
	)
	if err != nil {
		panic(err)
	}
	scene.CubeMap = cubeMap

	// Create BVH
	bvhRoot := utils.NewBVHNode(scene.World.Objects, 0, len(scene.World.Objects))
	scene.World = utils.HittableList{Objects: []utils.Hittable{bvhRoot}}

	world = scene.World
	cam = scene.Cam
	cam.Cube = *scene.CubeMap
	cam.SkipCube = scene.SkipBackground

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

package utils

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Camera struct {
	AspectRatio, pixelSamplesScale, Vfov                                float64
	ImageWidth, imageHeight, SamplesPerPixel, MaxDepth                  int
	center, pixel00Loc, pixelDeltaU, pixelDeltaV, LookFrom, LookAt, Vup Vec3
	u, v, w, defocusDiskU, defocusDiskV                                 Vec3
	DefocusAngle, Focusdist                                             float64
}
type Tile struct {
	x, y          int // Top-left corner
	width, height int
}

func (c *Camera) Render(world HittableList, display *DisplayBuffer) {
	c.initialize()
	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	t := time.Now()
	fmt.Fprintf(file, "P6\n%d %d\n%d\n", c.ImageWidth, c.imageHeight, 255)
	pixels := make([]byte, c.imageHeight*c.ImageWidth*3)

	tileWidth := 32
	tileHeight := 32
	numTilesX := (c.ImageWidth + tileWidth - 1) / tileWidth
	numTilesY := (c.imageHeight + tileHeight - 1) / tileHeight
	totalTiles := numTilesX * numTilesY

	tileChannel := make(chan Tile, totalTiles)
	resultChannel := make(chan struct {
		tile  Tile
		color []byte
	}, totalTiles)

	numWorkers := runtime.NumCPU() + 2
	var wg sync.WaitGroup
	var completedTiles atomic.Int32

	// Modified progress reporting goroutine with display updates
	go func() {
		for {
			completed := completedTiles.Load()
			if completed >= int32(totalTiles) || display.ShouldClose() {
				break
			}
			fmt.Printf("\033[1A\033[K")
			fmt.Printf("Progress: %.1f%% (%d/%d tiles)\n",
				float64(completed)/float64(totalTiles)*100,
				completed, totalTiles)

			// 60fps : )
			display.Refresh()
			time.Sleep(time.Second / 60)
		}
	}()

	// Worker function remains mostly the same
	worker := func(id int) {
		defer wg.Done()

		for tile := range tileChannel {
			effectiveWidth := min(tileWidth, c.ImageWidth-tile.x)
			effectiveHeight := min(tileHeight, c.imageHeight-tile.y)
			tileBuffer := make([]byte, effectiveWidth*effectiveHeight*3)

			for dy := 0; dy < effectiveHeight; dy++ {
				for dx := 0; dx < effectiveWidth; dx++ {
					if display.ShouldClose() {
						return
					}
					x := tile.x + dx
					y := tile.y + dy

					pixelColor := Vec3{0, 0, 0}
					for sample := 0; sample < c.SamplesPerPixel; sample++ {
						ray := c.getRay(x, y)
						pixelColor = pixelColor.PlusEq(rayColor(&ray, c.MaxDepth, &world))
					}

					finalColor := pixelColor.TimesConst(c.pixelSamplesScale)
					pixelIndex := (dy*effectiveWidth + dx) * 3
					WriteColor(tileBuffer, pixelIndex, finalColor)

					// Update display buffer
					intensity := Interval{0.000, 0.999}
					display.UpdatePixel(x, y, color.RGBA{
						R: uint8(int(256 * intensity.clamp(LinearToGamma(finalColor.X)))),
						G: uint8(int(256 * intensity.clamp(LinearToGamma(finalColor.Y)))),
						B: uint8(int(256 * intensity.clamp(LinearToGamma(finalColor.Z)))),
						A: 255,
					})
				}
			}

			resultChannel <- struct {
				tile  Tile
				color []byte
			}{tile, tileBuffer}

			completedTiles.Add(1)
		}
	}

	// Start workers
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go worker(w)
	}

	// Generate tiles
	go func() {
		for ty := 0; ty < c.imageHeight; ty += tileHeight {
			for tx := 0; tx < c.ImageWidth; tx += tileWidth {
				if display.ShouldClose() {
					close(tileChannel)
					return
				}
				tileChannel <- Tile{
					x:      tx,
					y:      ty,
					width:  min(tileWidth, c.ImageWidth-tx),
					height: min(tileHeight, c.imageHeight-ty),
				}
			}
		}
		close(tileChannel)
	}()

	// Collect results
	go func() {
		for result := range resultChannel {
			// Copy tile data to main pixel buffer
			for y := 0; y < result.tile.height; y++ {
				srcOffset := y * result.tile.width * 3
				dstOffset := ((result.tile.y+y)*c.ImageWidth + result.tile.x) * 3
				copy(pixels[dstOffset:], result.color[srcOffset:srcOffset+result.tile.width*3])
			}
		}
	}()

	// Wait for completion
	wg.Wait()
	close(resultChannel)

	// Only write the file if we didn't close the window
	if !display.ShouldClose() {
		_, err = file.Write(pixels)
		if err != nil {
			return
		}

		fmt.Printf("\033[1A\033[K")
		fmt.Printf("Done in: %v\n", time.Since(t))
		fmt.Printf("Image size: %d x %d\n", c.ImageWidth, c.imageHeight)

		//openFile("goimage.ppm")
	}
}

func (c *Camera) initialize() {
	c.imageHeight = int(float64(c.ImageWidth) / c.AspectRatio)
	if c.imageHeight < 0 {
		c.imageHeight = 1
	}

	c.pixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)
	c.center = c.LookFrom

	theta := DegreesToRadians(c.Vfov)
	h := math.Tan(theta / 2)
	viewportHeight := 2 * h * c.Focusdist
	viewportWidth := viewportHeight * (float64(c.ImageWidth) / float64(c.imageHeight))

	c.w = c.LookFrom.MinusEq(c.LookAt).UnitVector()
	c.u = c.Vup.Cross(c.w).UnitVector()
	c.v = c.w.Cross(c.u)

	viewportU := c.u.TimesConst(viewportWidth)
	viewportV := c.v.Neg().TimesConst(viewportHeight)

	c.pixelDeltaU = viewportU.TimesConst(1.0 / float64(c.ImageWidth))
	c.pixelDeltaV = viewportV.TimesConst(1.0 / float64(c.imageHeight))

	viewportUpperLeft := c.center.MinusEq(c.w.TimesConst(c.Focusdist)).MinusEq(viewportU.TimesConst(0.5)).MinusEq(viewportV.TimesConst(0.5))
	c.pixel00Loc = c.pixelDeltaU.PlusEq(c.pixelDeltaV).TimesConst(0.5).PlusEq(viewportUpperLeft)

	defocusRadius := c.Focusdist * math.Tan(DegreesToRadians(c.DefocusAngle/2))
	c.defocusDiskU = c.u.TimesConst(defocusRadius)
	c.defocusDiskV = c.v.TimesConst(defocusRadius)
}
func rayColor(r *Ray, depth int, world Hittable) Vec3 {
	if depth <= 0 {
		return Vec3{0, 0, 0}
	}

	var rec HitRecord

	if world.Hit(r, Interval{0.001, math.Inf(+1)}, &rec) {
		var scattered Ray
		var attenuation Vec3
		if rec.Mat.Scatter(r, &scattered, &attenuation, &rec) {
			return rayColor(&scattered, depth-1, world).TimesEq(attenuation)
		}
		return Vec3{0, 0, 0}
	}

	unitDirection := r.Direction.Normalize()
	a := 0.5 * (unitDirection.Y + 1.0)
	white := Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue := Vec3{X: 0.5, Y: 0.7, Z: 1.0}
	return white.TimesConst(1.0 - a).PlusEq(blue.TimesConst(a))
}

// auto ray_origin = (defocus_angle <= 0) ? center : defocus_disk_sample();
func (c *Camera) getRay(i, j int) Ray {
	offset := sampleSquare()
	pixelSample := c.pixel00Loc.PlusEq(c.pixelDeltaU.TimesConst(float64(i) + offset.X)).PlusEq(c.pixelDeltaV.TimesConst(float64(j) + offset.Y))
	var rayOrigin Vec3
	if c.DefocusAngle <= 0 {
		rayOrigin = c.center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.MinusEq(rayOrigin)

	return Ray{rayOrigin, rayDirection}
}
func sampleSquare() Vec3 {
	return Vec3{RandomFloat() - 0.5, RandomFloat() - 0.5, 0}
}
func (c *Camera) defocusDiskSample() Vec3 {
	p := RandomInUnitDisk()
	return c.center.PlusEq(c.defocusDiskU.TimesConst(p.X)).PlusEq(c.defocusDiskV.TimesConst(p.Y))
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

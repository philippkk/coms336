package utils

import (
	"fmt"
	"image/color"
	"math"
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

func (c *Camera) Render(world HittableList, display *DisplayBuffer, pixels []byte) time.Duration {
	c.initialize()
	t := time.Now()
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

	//displayWidth, displayHeight := display.Win.Bounds().Max.XY()

	go func() {
		for {
			completed := completedTiles.Load()
			if completed >= int32(totalTiles) || display.ShouldClose() {
				break
			}

			//fmt.Printf("\033[1A\033[K")
			//fmt.Printf("Progress: %.1f%% (%d/%d tiles)\n",
			//	float64(completed)/float64(totalTiles)*100,
			//	completed, totalTiles)

			// 60fps
			display.Refresh()
			time.Sleep(time.Second / 60)
		}
	}()

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

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go worker(w)
	}

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

	go func() {
		for result := range resultChannel {
			for y := 0; y < result.tile.height; y++ {
				srcOffset := y * result.tile.width * 3
				dstOffset := ((result.tile.y+y)*c.ImageWidth + result.tile.x) * 3
				copy(pixels[dstOffset:], result.color[srcOffset:srcOffset+result.tile.width*3])
			}
		}
	}()

	wg.Wait()
	time.Sleep(time.Second / 2)
	close(resultChannel)
	fmt.Printf("\033[1A\033[K")
	fmt.Printf("Done in: %v\n", time.Since(t))
	return time.Since(t)
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
	rayTime := RandomFloat()

	return Ray{rayOrigin, rayDirection, rayTime}
}
func sampleSquare() Vec3 {
	return Vec3{RandomFloat() - 0.5, RandomFloat() - 0.5, 0}
}
func (c *Camera) defocusDiskSample() Vec3 {
	p := RandomInUnitDisk()
	return c.center.PlusEq(c.defocusDiskU.TimesConst(p.X)).PlusEq(c.defocusDiskV.TimesConst(p.Y))
}

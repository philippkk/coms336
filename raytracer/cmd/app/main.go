package main

import (
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	width, height := 2000, 2000
	maxColorValue := 255

	file, err := os.Create("goimage.ppm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "P6\n%d %d\n%d\n", width, height, maxColorValue)
	
	var pixels []byte
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			color := utils.NewVec3(float64(j)/float64(width-1),
				float64(i)/float64(height-1),
				255)
			utils.WriteColor(&pixels, *color)
		}
	}

	_, err = file.Write(pixels)
	if err != nil {
		return
	}

	fmt.Println("Done.")

	openFile("goimage.ppm")
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

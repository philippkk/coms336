package model

import (
	"bufio"
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"strconv"
	"strings"
)

type MTLMaterial struct {
	Name string

	// Optical properties
	AmbientColor  mgl32.Vec3
	DiffuseColor  mgl32.Vec3
	SpecularColor mgl32.Vec3

	// Texture maps
	AmbientTexture  string
	DiffuseTexture  string
	SpecularTexture string

	// Other properties
	Shininess         float32
	Opacity           float32
	IlluminationModel int
}

func ParseMTLFile(filename string) map[string]MTLMaterial {
	materials := make(map[string]MTLMaterial)
	currentMaterial := MTLMaterial{}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "newmtl":
			// Save previous material if exists
			if currentMaterial.Name != "" {
				materials[currentMaterial.Name] = currentMaterial
			}
			// Start new material
			currentMaterial = MTLMaterial{
				Name: fields[1],
				// Set default values
				AmbientColor:  mgl32.Vec3{0.2, 0.2, 0.2},
				DiffuseColor:  mgl32.Vec3{0.8, 0.8, 0.8},
				SpecularColor: mgl32.Vec3{1.0, 1.0, 1.0},
				Opacity:       1.0,
				Shininess:     0.0,
			}
		case "Ka":
			// Ambient color
			currentMaterial.AmbientColor = parseVec3(fields[1:])
		case "Kd":
			// Diffuse color
			currentMaterial.DiffuseColor = parseVec3(fields[1:])
		case "Ks":
			// Specular color
			currentMaterial.SpecularColor = parseVec3(fields[1:])
		case "Ns":
			// Shininess
			shininess, _ := strconv.ParseFloat(fields[1], 32)
			currentMaterial.Shininess = float32(shininess)
		case "d":
			// Opacity
			opacity, _ := strconv.ParseFloat(fields[1], 32)
			currentMaterial.Opacity = float32(opacity)
		case "illum":
			// Illumination model
			model, _ := strconv.Atoi(fields[1])
			currentMaterial.IlluminationModel = model
		case "map_Ka":
			// Ambient texture
			currentMaterial.AmbientTexture = fields[1]
		case "map_Kd":
			// Diffuse texture
			currentMaterial.DiffuseTexture = fields[1]
		case "map_Ks":
			// Specular texture
			currentMaterial.SpecularTexture = fields[1]
		}
	}

	// Add last material
	if currentMaterial.Name != "" {
		materials[currentMaterial.Name] = currentMaterial
	}

	return materials
}

// Helper function to parse Vec3
func parseVec3(values []string) mgl32.Vec3 {
	var vec mgl32.Vec3
	for i := 0; i < 3 && i < len(values); i++ {
		val, _ := strconv.ParseFloat(values[i], 32)
		vec[i] = float32(val)
	}
	return vec
}

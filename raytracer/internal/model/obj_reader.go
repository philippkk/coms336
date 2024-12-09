package model

import (
	"bufio"
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"github.com/philippkk/coms336/raytracer/internal/utils/material"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// Model is a renderable collection of vecs.
type Model struct {
	// For the v, vt and vn in the obj file.
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2

	// For the fun "f" in the obj file.
	VecIndices, NormalIndices, UvIndices []float32

	MaterialLib          map[string]MTLMaterial
	MaterialIndexChanges []int    // Indices where material changes occur
	MaterialNames        []string // Names of materials at each change point
}

func NewModel(obj, mtlFile string) Model {
	objFile, err := os.Open(obj)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	model := Model{
		MaterialLib:          ParseMTLFile(mtlFile),
		MaterialIndexChanges: []int{},
		MaterialNames:        []string{},
	}

	model.MaterialLib = ParseMTLFile(mtlFile)

	// Scan the file line by line
	scanner := bufio.NewScanner(objFile)
	triangleCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "v":
			// Vertex positions
			var vec mgl32.Vec3
			fmt.Sscanf(fields[1], "%f", &vec[0])
			fmt.Sscanf(fields[2], "%f", &vec[1])
			fmt.Sscanf(fields[3], "%f", &vec[2])
			model.Vecs = append(model.Vecs, vec)
		case "vn":
			// Vertex normals
			var vec mgl32.Vec3
			fmt.Sscanf(fields[1], "%f", &vec[0])
			fmt.Sscanf(fields[2], "%f", &vec[1])
			fmt.Sscanf(fields[3], "%f", &vec[2])
			model.Normals = append(model.Normals, vec)
		case "vt":
			// Texture coordinates
			var uv mgl32.Vec2
			fmt.Sscanf(fields[1], "%f", &uv[0])
			fmt.Sscanf(fields[2], "%f", &uv[1])
			model.Uvs = append(model.Uvs, uv)
		case "f":
			// Faces (triangulate)
			vertices := fields[1:]
			for i := 1; i+1 < len(vertices); i++ {
				processFace(vertices[0], &model)
				processFace(vertices[i], &model)
				processFace(vertices[i+1], &model)
				triangleCount++
			}
		case "usemtl":
			// Set current material
			model.MaterialIndexChanges = append(model.MaterialIndexChanges, triangleCount)
			model.MaterialNames = append(model.MaterialNames, fields[1])
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return model
}

func processFace(face string, model *Model) {
	var vIdx, tIdx, nIdx int
	segments := strings.Split(face, "/")

	if len(segments) > 0 && segments[0] != "" {
		vIdx, _ = strconv.Atoi(segments[0])
	}
	if len(segments) > 1 && segments[1] != "" {
		tIdx, _ = strconv.Atoi(segments[1])
	}
	if len(segments) > 2 && segments[2] != "" {
		nIdx, _ = strconv.Atoi(segments[2])
	}

	// OBJ indices are 1-based, convert to 0-based
	if vIdx > 0 {
		model.VecIndices = append(model.VecIndices, float32(vIdx-1))
	}
	if tIdx > 0 {
		model.UvIndices = append(model.UvIndices, float32(tIdx-1))
	}
	if nIdx > 0 {
		model.NormalIndices = append(model.NormalIndices, float32(nIdx-1))
	}
}

func (model Model) ToTriangles(defaultMat utils.Material, name string) []objects.Triangle {
	var triangles []objects.Triangle

	// Pre-cache materials to avoid repeated texture loading
	materialCache := make(map[string]utils.Material)
	materialCache["default"] = defaultMat

	// Pre-cache materials for better performance
	for materialName, mtlMaterial := range model.MaterialLib {
		if mtlMaterial.DiffuseTexture != "" {
			// Debug print for the specific material
			fmt.Printf("Loading material: %s with texture: %s\n", materialName, mtlMaterial.DiffuseTexture)

			texture, err := utils.NewImageTexture("internal/model/" + name + "/" + mtlMaterial.DiffuseTexture)
			if err == nil {
				// Print texture details after successful loading
				fmt.Printf("Successfully loaded texture for %s: %dx%d\n",
					materialName, texture.Width, texture.Height)

				if mtlMaterial.DiffuseTexture == "1_Diffuse.png" {
					materialCache[materialName] = material.Metal{utils.Vec3{0, 0.16078431372, 0.51764705882}, 0.4}
				} else {
					materialCache[materialName] = material.NewLambertian(texture)
				}
			}
		}
	}

	// Pre-allocate the slice to reduce allocations
	triangles = make([]objects.Triangle, 0, len(model.VecIndices)/3)

	// Use a single progress tracker
	totalTriangles := len(model.VecIndices) / 3
	lastProgressPrint := time.Now()

	for i := 0; i < len(model.VecIndices); i += 3 {
		// Get vertex positions (existing code)
		v0Index := int(model.VecIndices[i])
		v1Index := int(model.VecIndices[i+1])
		v2Index := int(model.VecIndices[i+2])

		// Get vertex positions
		v0 := model.Vecs[v0Index]
		v1 := model.Vecs[v1Index]
		v2 := model.Vecs[v2Index]

		t0 := utils.Vec3{X: float64(v0.X()), Y: float64(v0.Y()), Z: float64(v0.Z())}
		t1 := utils.Vec3{X: float64(v1.X()), Y: float64(v1.Y()), Z: float64(v1.Z())}
		t2 := utils.Vec3{X: float64(v2.X()), Y: float64(v2.Y()), Z: float64(v2.Z())}

		// Get vertex normals
		var n0, n1, n2 utils.Vec3
		if len(model.NormalIndices) > i+2 {
			nIndex0 := int(model.NormalIndices[i])
			nIndex1 := int(model.NormalIndices[i+1])
			nIndex2 := int(model.NormalIndices[i+2])

			if nIndex0 < len(model.Normals) && nIndex1 < len(model.Normals) && nIndex2 < len(model.Normals) {
				vn0 := model.Normals[nIndex0]
				vn1 := model.Normals[nIndex1]
				vn2 := model.Normals[nIndex2]

				n0 = utils.Vec3{X: float64(vn0.X()), Y: float64(vn0.Y()), Z: float64(vn0.Z())}
				n1 = utils.Vec3{X: float64(vn1.X()), Y: float64(vn1.Y()), Z: float64(vn1.Z())}
				n2 = utils.Vec3{X: float64(vn2.X()), Y: float64(vn2.Y()), Z: float64(vn2.Z())}
			}
		}

		// Find the material for this triangle (existing code)
		triangleMat := defaultMat
		if len(model.MaterialIndexChanges) > 0 {
			materialIndex := findMaterialIndexForTriangle(i, model.MaterialIndexChanges)
			if materialIndex != -1 {
				materialName := model.MaterialNames[materialIndex]
				if cachedMat, exists := materialCache[materialName]; exists {
					triangleMat = cachedMat
				}
			}
		}

		// Create triangle with normals
		var triangle objects.Triangle
		if len(model.NormalIndices) > 0 {
			triangle = objects.CreateTriangleWithNormals(
				t0, t1, t2,
				n0, n1, n2,
				triangleMat,
				getUVCoord(model, i),
				getUVCoord(model, i+1),
				getUVCoord(model, i+2),
			)
		} else {
			triangle = objects.CreateTriangleWithUV(
				t0, t1, t2,
				triangleMat,
				getUVCoord(model, i),
				getUVCoord(model, i+1),
				getUVCoord(model, i+2),
			)
		}

		triangles = append(triangles, triangle)
		// Efficient progress tracking
		if time.Since(lastProgressPrint) > 500*time.Millisecond {
			progress := float64(i/3) / float64(totalTriangles) * 100
			fmt.Printf("\rProcessing triangles: %.2f%% (%d/%d)", progress, i/3, totalTriangles)
			lastProgressPrint = time.Now()
		}
	}

	fmt.Println("\rTriangle processing complete.                    ")
	return triangles
}

// Optimized material index finding
func findMaterialIndexForTriangle(triangleIndex int, materialChanges []int) int {
	// Binary search would be even faster for large models
	for i := len(materialChanges) - 1; i >= 0; i-- {
		if materialChanges[i] <= triangleIndex/3 {
			return i
		}
	}
	return -1
}

// Helper function to get UV coordinates
func wrapUV(value float64) float64 {
	wrapped := value - math.Floor(value) // This handles both positive and negative values
	if wrapped < 0 {
		wrapped += 1.0
	}
	return wrapped
}

func getUVCoord(model Model, index int) utils.Vec2 {
	if len(model.UvIndices) > index && int(model.UvIndices[index]) < len(model.Uvs) {
		uv := model.Uvs[int(model.UvIndices[index])]
		return utils.Vec2{
			X: wrapUV(float64(uv.X())),
			Y: wrapUV(float64(uv.Y())),
		}
	}
	return utils.Vec2{X: 0, Y: 0}
}

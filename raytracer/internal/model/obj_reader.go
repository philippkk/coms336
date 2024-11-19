package model

import (
	"bufio"
	"fmt"
	"github.com/philippkk/coms336/raytracer/internal/objects"
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

// Model is a renderable collection of vecs.
type Model struct {
	// For the v, vt and vn in the obj file.
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2

	// For the fun "f" in the obj file.
	VecIndices, NormalIndices, UvIndices []float32
}

func NewModel(file string) Model {
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	model := Model{}

	// Scan the file line by line
	scanner := bufio.NewScanner(objFile)
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
			}
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

	if vIdx != 0 {
		model.VecIndices = append(model.VecIndices, float32(vIdx-1))
	}
	if tIdx != 0 {
		model.UvIndices = append(model.UvIndices, float32(tIdx-1))
	}
	if nIdx != 0 {
		model.NormalIndices = append(model.NormalIndices, float32(nIdx-1))
	}
}

func (model Model) ToTriangles(mat utils.Material) []objects.Triangle {
	var triangles []objects.Triangle

	for i := 0; i < len(model.VecIndices); i += 3 {
		v0Index := int(model.VecIndices[i])
		v1Index := int(model.VecIndices[i+1])
		v2Index := int(model.VecIndices[i+2])

		if v0Index < 0 || v1Index < 0 || v2Index < 0 ||
			v0Index >= len(model.Vecs) || v1Index >= len(model.Vecs) || v2Index >= len(model.Vecs) {
			fmt.Println("Skipping invalid triangle indices:", v0Index, v1Index, v2Index)
			continue
		}

		v0 := model.Vecs[v0Index]
		v1 := model.Vecs[v1Index]
		v2 := model.Vecs[v2Index]

		t0 := utils.Vec3{X: float64(v0.X()), Y: float64(v0.Y()), Z: float64(v0.Z())}
		t1 := utils.Vec3{X: float64(v1.X()), Y: float64(v1.Y()), Z: float64(v1.Z())}
		t2 := utils.Vec3{X: float64(v2.X()), Y: float64(v2.Y()), Z: float64(v2.Z())}

		triangles = append(triangles, objects.CreateTriangle(t0, t1, t2, mat))
	}

	return triangles
}

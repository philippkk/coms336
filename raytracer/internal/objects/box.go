package objects

import (
	"github.com/philippkk/coms336/raytracer/internal/utils"
	"math"
)

func CreateBox(a, b utils.Vec3, mat utils.Material) []utils.Hittable {
	sides := make([]utils.Hittable, 0, 6)

	// Get min and max points
	min := utils.Vec3{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
		Z: math.Min(a.Z, b.Z),
	}
	max := utils.Vec3{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
		Z: math.Max(a.Z, b.Z),
	}

	// Create the edge vectors
	dx := utils.Vec3{X: max.X - min.X}
	dy := utils.Vec3{Y: max.Y - min.Y}
	dz := utils.Vec3{Z: max.Z - min.Z}

	// Front face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: min.X, Y: min.Y, Z: max.Z},
		dx,
		dy,
		mat,
	))

	// Right face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: max.X, Y: min.Y, Z: max.Z},
		dz.Neg(),
		dy,
		mat,
	))

	// Back face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: max.X, Y: min.Y, Z: min.Z},
		dx.Neg(),
		dy,
		mat,
	))

	// Left face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: min.X, Y: min.Y, Z: min.Z},
		dz,
		dy,
		mat,
	))

	// Top face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: min.X, Y: max.Y, Z: max.Z},
		dx,
		dz.Neg(),
		mat,
	))

	// Bottom face
	sides = append(sides, CreateQuad(
		utils.Vec3{X: min.X, Y: min.Y, Z: min.Z},
		dx,
		dz,
		mat,
	))

	return sides
}

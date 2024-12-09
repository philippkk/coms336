package utils

import (
	"math"
)

// Face identifiers for the cube map
const (
	RIGHT = iota
	LEFT
	TOP
	BOTTOM
	FRONT
	BACK
)

type CubeMap struct {
	faces [6]*ImageTexture
}

// NewCubeMap creates a cube map from 6 image files
func NewCubeMap(right, left, top, bottom, front, back string) (*CubeMap, error) {
	cm := &CubeMap{}

	// Load all face textures
	files := []string{right, left, top, bottom, front, back}
	for i, file := range files {
		texture, err := NewImageTexture(file)
		if err != nil {
			return nil, err
		}
		cm.faces[i] = texture
	}

	return cm, nil
}

func (cm *CubeMap) SampleCubeMap(dir Vec3) Vec3 {
	// Find which face to use and convert to UV coordinates
	face, u, v := cm.directionToUV(dir)

	// Sample the appropriate face texture
	return cm.faces[face].Value(u, v, Vec3{})
}

// directionToUV converts a direction vector to face index and UV coordinates
func (cm *CubeMap) directionToUV(dir Vec3) (face int, u, v float64) {
	// Find the largest component to determine which face to use
	absX := math.Abs(dir.X)
	absY := math.Abs(dir.Y)
	absZ := math.Abs(dir.Z)
	maxAxis := math.Max(absX, math.Max(absY, absZ))

	// Calculate UV based on which face we're using
	switch maxAxis {
	case absX:
		if dir.X > 0 { // Right face
			face = RIGHT
			u = (-dir.Z/absX + 1.0) / 2.0
			v = (dir.Y/absX + 1.0) / 2.0 // Removed negative
		} else { // Left face
			face = LEFT
			u = (dir.Z/absX + 1.0) / 2.0
			v = (dir.Y/absX + 1.0) / 2.0 // Removed negative
		}
	case absY:
		if dir.Y > 0 { // Top face
			face = TOP
			u = (dir.X/absY + 1.0) / 2.0
			v = (dir.Z/absY + 1.0) / 2.0
		} else { // Bottom face
			face = BOTTOM
			u = (dir.X/absY + 1.0) / 2.0
			v = (-dir.Z/absY + 1.0) / 2.0
		}
	case absZ:
		if dir.Z > 0 { // Front face
			face = FRONT
			u = (dir.X/absZ + 1.0) / 2.0
			v = (dir.Y/absZ + 1.0) / 2.0 // Removed negative
		} else { // Back face
			face = BACK
			u = (-dir.X/absZ + 1.0) / 2.0
			v = (dir.Y/absZ + 1.0) / 2.0 // Removed negative
		}
	}

	// Ensure UV coordinates are in [0,1] range
	u = math.Max(0, math.Min(1, u))
	v = math.Max(0, math.Min(1, v))

	return face, u, v
}

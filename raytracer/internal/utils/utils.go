package utils

import "math/rand"

const (
	pi = 3.1415926535897932385
)

func DegreesToRadians(degrees float64) float64 {
	return degrees * pi / 180.0
}

func RandomFloat() float64 {
	return rand.Float64() // rand_max + 1???
}

func RandomFloatInRange(min, max float64) float64 {
	return min + (max-min)*RandomFloat()
}

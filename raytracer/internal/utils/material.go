package utils

type Material interface {
	Scatter(rIn, scattered *Ray, attenuation *Vec3, rec HitRecord) bool
}

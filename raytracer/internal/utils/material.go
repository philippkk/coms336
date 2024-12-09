package utils

type Material interface {
	Scatter(rIn, scattered *Ray, attenuation *Vec3, rec *HitRecord) bool

	ColorEmitted(u, v float64, p Vec3) Vec3
}

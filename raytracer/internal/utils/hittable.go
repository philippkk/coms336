package utils

type HitRecord struct {
	P, Normal Vec3
	T         float64
}

type Hittable interface {
	hit(ray *Ray, rayTmin, rayTmax float64, rec HitRecord) bool
	test() int
}

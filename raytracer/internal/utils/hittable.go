package utils

type HitRecord struct {
	P, Normal Vec3
	T         float64
	FrontFace bool
	Mat       Material
}

func (h *HitRecord) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	h.FrontFace = r.Direction.Dot(outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Neg()
	}
}

type Hittable interface {
	BoundingBox() AABB
	Hit(ray *Ray, rayT Interval, rec *HitRecord) bool
}

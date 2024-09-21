package utils

type HitRecord struct {
	P, Normal Vec3
	T         float64
	FrontFace bool
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
	Hit(ray *Ray, rayTmin, rayTmax float64, rec *HitRecord) bool
}

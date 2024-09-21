package utils

type HittableList struct {
	objects []Hittable
}

func (h *HittableList) Add(object Hittable) {
	h.objects = append(h.objects, object)
}

func (h *HittableList) Clear() {
	h.objects = nil

}

func (h *HittableList) Hit(ray *Ray, rayTmin, rayTmax float64, rec HitRecord) bool {
	var tempRec HitRecord
	var hitAnything bool
	closestSoFar := rayTmax

	for _, obj := range h.objects {
		if obj.Hit(ray, rayTmin, closestSoFar, tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			rec = tempRec
		}
	}

	return hitAnything
}

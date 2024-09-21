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

func (h *HittableList) Hit(ray *Ray, rayT Interval, rec *HitRecord) bool {
	var tempRec HitRecord
	var hitAnything bool
	closestSoFar := rayT.Max

	for _, obj := range h.objects {
		if obj.Hit(ray, Interval{rayT.Min, closestSoFar}, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}

	return hitAnything
}

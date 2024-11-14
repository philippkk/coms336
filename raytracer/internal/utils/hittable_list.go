package utils

type HittableList struct {
	Objects []Hittable
	Box     AABB
}

func (h *HittableList) GetSize() int {
	return len(h.Objects)
}

func (h *HittableList) BoundingBox() AABB {
	return h.Box
}

func (h *HittableList) Add(object Hittable) {
	h.Objects = append(h.Objects, object)
	h.Box = SurroundingBox(h.Box, object.BoundingBox())
}
func (h *HittableList) Clear() {
	h.Objects = nil

}
func (h *HittableList) Hit(ray *Ray, rayT Interval, rec *HitRecord) bool {
	var tempRec HitRecord
	var hitAnything bool
	closestSoFar := rayT.Max

	for _, obj := range h.Objects {
		if obj.Hit(ray, Interval{rayT.Min, closestSoFar}, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}

	return hitAnything
}

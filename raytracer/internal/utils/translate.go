package utils

type Translate struct {
	Offset Vec3
	Object Hittable
}

func (t *Translate) Hit(r *Ray, rayT Interval, rec *HitRecord) bool {
	offsetR := Ray{r.Origin.MinusEq(t.Offset), r.Direction, r.Tm}

	if !t.Object.Hit(&offsetR, rayT, rec) {
		return false
	}
	rec.P = rec.P.PlusEq(t.Offset)
	return true
}

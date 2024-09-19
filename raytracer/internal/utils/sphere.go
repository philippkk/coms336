package utils

import "fmt"

type Sphere struct {
	Val int
}

func (s Sphere) hit(ray *Ray, rayTmin, rayTmax float64, rec HitRecord) bool {
	//TODO implement me
	panic("implement me")
}

func (s Sphere) test() int {
	return s.Val
}

type Tri struct {
	Val int
}

func (s Tri) hit(ray *Ray, rayTmin, rayTmax float64, rec HitRecord) bool {
	//TODO implement me
	panic("implement me")
}

func (s Tri) test() int {
	return s.Val
}

func Measure(g Hittable) {
	fmt.Println(g.test())
}

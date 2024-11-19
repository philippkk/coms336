package utils

import (
	"math/rand/v2"
	"sort"
)

// BVHNode struct definition
type BVHNode struct {
	Left  Hittable
	Right Hittable
	Box   AABB
}

func (B BVHNode) BoundingBox() AABB {
	return B.Box
}

func (B BVHNode) Hit(ray *Ray, rayT Interval, rec *HitRecord) bool {
	if !B.Box.Hit(ray, rayT) {
		return false
	}
	hitL := B.Left.Hit(ray, rayT, rec)
	var hitR bool
	if hitL {
		hitR = B.Right.Hit(ray, Interval{rayT.Min, rec.T}, rec)
	} else {
		hitR = B.Right.Hit(ray, Interval{rayT.Min, rayT.Max}, rec)
	}
	return hitL || hitR
}

// Box comparison functions
func boxCompare(a, b Hittable, axisIndex int) bool {
	aBox := a.BoundingBox()
	bBox := b.BoundingBox()

	aInterval := aBox.AxisInterval(axisIndex)
	bInterval := bBox.AxisInterval(axisIndex)

	return aInterval.Min < bInterval.Min
}

func boxXCompare(a, b Hittable) bool {
	return boxCompare(a, b, 0)
}

func boxYCompare(a, b Hittable) bool {
	return boxCompare(a, b, 1)
}

func boxZCompare(a, b Hittable) bool {
	return boxCompare(a, b, 2)
}

// Constructor function
func NewBVHNode(objects []Hittable, start, end int) BVHNode {
	node := BVHNode{}

	axis := rand.IntN(3) // Choose between 0, 1, or 2

	objectSpan := end - start

	switch objectSpan {
	case 1:
		node.Left = objects[start]
		node.Right = objects[start]
	case 2:
		// Compare the two objects directly along the chosen axis
		// and arrange them in sorted order
		if boxCompare(objects[start], objects[start+1], axis) {
			node.Left = objects[start]
			node.Right = objects[start+1]
		} else {
			node.Left = objects[start+1]
			node.Right = objects[start]
		}
	default:
		// Sort objects along the chosen axis
		sort.Slice(objects[start:end], func(i, j int) bool {
			return boxCompare(objects[start+i], objects[start+j], axis)
		})

		mid := start + objectSpan/2
		node.Left = NewBVHNode(objects, start, mid)
		node.Right = NewBVHNode(objects, mid, end)
	}

	node.Box = SurroundingBox(node.Left.BoundingBox(), node.Right.BoundingBox())
	return node
}

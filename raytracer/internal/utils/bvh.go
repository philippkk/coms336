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

	// Randomly choose an axis to sort on
	axis := rand.IntN(2) + 1

	// Create comparison function based on axis
	comparator := func(i, j int) bool {
		if axis == 0 {
			return boxXCompare(objects[i], objects[j])
		} else if axis == 1 {
			return boxYCompare(objects[i], objects[j])
		}
		return boxZCompare(objects[i], objects[j])
	}

	objectSpan := end - start

	if objectSpan == 1 {
		node.Left = objects[start]
		node.Right = objects[start]
	} else if objectSpan == 2 {
		node.Left = objects[start]
		node.Right = objects[start+1]
	} else {
		// Sort the slice of objects
		sort.Slice(objects[start:end], func(i, j int) bool {
			return comparator(start+i, start+j)
		})

		mid := start + objectSpan/2
		node.Left = NewBVHNode(objects, start, mid)
		node.Right = NewBVHNode(objects, mid, end)
	}

	// Get the bounding box
	leftBox := node.Left.BoundingBox()
	rightBox := node.Right.BoundingBox()
	node.Box = SurroundingBox(leftBox, rightBox)

	return node
}

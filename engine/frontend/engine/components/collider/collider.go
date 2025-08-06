package collider

type Vec3 struct{ X, Y, Z float32 }

type Intersection interface {
	// point nearest to shape A center
	PointOnA() Vec3
	// point nearest to shape B center
	PointOnB() Vec3
}

type intersection struct {
	pointOnA, pointOnB Vec3
}

func (i *intersection) PointOnA() Vec3 { return i.pointOnA }
func (i *intersection) PointOnB() Vec3 { return i.pointOnB }

func NewIntersection(pointOnA, pointOnB Vec3) Intersection {
	return &intersection{
		pointOnA: pointOnA,
		pointOnB: pointOnB,
	}
}

//

// type Shape any

// shape should have relative position to transform.Position
// shape should have relative rotation to transform.Rotation
// shape should have relative size to tranform.Size (in percantages)
type Shape interface {
	Intersects(other Shape) (Intersection, error)
}

type Collider struct {
	Shape Shape
}

package collider

import (
	"frontend/engine/components/transform"
)

type box struct {
	t transform.Transform
}

func (shape *box) Intersects(other Shape) (Intersection, error) { return collides(shape, other) }

//

type sphere struct {
	t transform.Transform
}

func (shape *sphere) Intersects(other Shape) (Intersection, error) { return collides(shape, other) }

//

type ray struct {
	t transform.Transform
}

func (shape *ray) Intersects(other Shape) (Intersection, error) { return collides(shape, other) }

//

//

func NewBox(t transform.Transform) Shape    { return &box{t: t} }
func NewRect(t transform.Transform) Shape   { return &box{t: t} }
func NewSphere(t transform.Transform) Shape { return &sphere{t: t} }
func NewCircle(t transform.Transform) Shape { return &sphere{t: t} }
func NewRay(t transform.Transform) Shape    { return &ray{t: t} }

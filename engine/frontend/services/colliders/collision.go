package colliders

type Collision interface {
	Intersections() []Intersection
	Reverse() Collision
}

type collision struct {
	intersections []Intersection
}

func (collision *collision) Intersections() []Intersection { return collision.intersections }
func (collision *collision) Reverse() Collision {
	intersections := make([]Intersection, 0, len(collision.intersections))
	for _, intersection := range collision.intersections {
		intersections = append(intersections, intersection.Reverse())
	}
	return NewCollision(intersections...)
}

func NewCollision(intersections ...Intersection) Collision {
	return &collision{intersections: intersections}
}

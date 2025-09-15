package anchor

import (
	"frontend/engine/components/transform"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type RelativeChange uint8

const (
	Ignore RelativeChange = iota
	// ChangeParent
	// ChangeChild
)

type ParentAnchor struct {
	Parent            ecs.EntityID
	OnChildChange     RelativeChange
	RelativeTransform transform.Transform
	// locks refer to object center
	// every lock axis should be between 0 and 1
	ParentPivot transform.PivotPoint
}

func NewParentAnchor(parent ecs.EntityID) ParentAnchor {
	return ParentAnchor{
		parent,
		Ignore,
		transform.NewTransform(),
		transform.NewPivotPoint(mgl32.Vec3{.5, .5, .5}),
	}
}

func (c ParentAnchor) Ptr() *ParentAnchor { return &c }
func (c *ParentAnchor) Val() ParentAnchor { return *c }

func (c *ParentAnchor) SetParent(entity ecs.EntityID) *ParentAnchor {
	c.Parent = entity
	return c
}

func (c *ParentAnchor) SetChildChange(change RelativeChange) *ParentAnchor {
	c.OnChildChange = change
	return c
}

func (c *ParentAnchor) SetRelativeTransform(transform transform.Transform) *ParentAnchor {
	c.RelativeTransform = transform
	return c
}

func (c *ParentAnchor) SetPivotPoint(pivot mgl32.Vec3) *ParentAnchor {
	c.ParentPivot = transform.NewPivotPoint(pivot)
	return c
}

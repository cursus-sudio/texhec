package anchor

import (
	"frontend/engine/transform"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type RelativeChange uint8

const (
	Ignore RelativeChange = iota
	// ChangeParent
	// ChangeChild
)

type ParentAnchorComponent struct {
	Parent            ecs.EntityID
	OnChildChange     RelativeChange
	RelativeTransform transform.TransformComponent
	// locks refer to object center
	// every lock axis should be between 0 and 1
	ParentPivot transform.PivotPointComponent
	// offset is a tool to create margin
	Offset mgl32.Vec3
}

func NewParentAnchor(parent ecs.EntityID) ParentAnchorComponent {
	return ParentAnchorComponent{
		parent,
		Ignore,
		transform.NewTransform(),
		transform.NewPivotPoint(mgl32.Vec3{.5, .5, .5}),
		mgl32.Vec3{},
	}
}

func (c ParentAnchorComponent) Ptr() *ParentAnchorComponent { return &c }
func (c *ParentAnchorComponent) Val() ParentAnchorComponent { return *c }

func (c *ParentAnchorComponent) SetParent(entity ecs.EntityID) *ParentAnchorComponent {
	c.Parent = entity
	return c
}

func (c *ParentAnchorComponent) SetChildChange(change RelativeChange) *ParentAnchorComponent {
	c.OnChildChange = change
	return c
}

func (c *ParentAnchorComponent) SetRelativeTransform(transform transform.TransformComponent) *ParentAnchorComponent {
	c.RelativeTransform = transform
	return c
}

func (c *ParentAnchorComponent) SetPivotPoint(pivot mgl32.Vec3) *ParentAnchorComponent {
	c.ParentPivot = transform.NewPivotPoint(pivot)
	return c
}

func (c *ParentAnchorComponent) SetOffset(offset mgl32.Vec3) *ParentAnchorComponent {
	c.Offset = offset
	return c
}

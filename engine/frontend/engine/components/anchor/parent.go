package anchor

import (
	"frontend/engine/components/transform"
	"frontend/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type RelativeChange uint8

const (
	Ignore RelativeChange = iota
	// ChangeParent
	// ChangeChild
)

type parentAnchor struct {
	Parent            ecs.EntityID
	OnChildChange     RelativeChange
	RelativeTransform transform.Transform
	// locks refer to object center
	// every lock axis should be between 0 and 1
	ParentPivot transform.PivotPoint
}

type ParentAnchor struct{ *parentAnchor }

func NewParentAnchor(parent ecs.EntityID) ParentAnchor {
	return ParentAnchor{&parentAnchor{
		parent,
		Ignore,
		transform.NewTransform(),
		transform.NewPivotPoint(mgl32.Vec3{.5, .5, .5}),
	}}
}

func (c ParentAnchor) SetParent(entity ecs.EntityID) ParentAnchor {
	return ParentAnchor{&parentAnchor{entity, c.OnChildChange, c.RelativeTransform, c.ParentPivot}}
}

func (c ParentAnchor) SetChildChange(change RelativeChange) ParentAnchor {
	return ParentAnchor{&parentAnchor{c.Parent, change, c.RelativeTransform, c.ParentPivot}}
}

func (c ParentAnchor) SetRelativeTransform(transform transform.Transform) ParentAnchor {
	return ParentAnchor{&parentAnchor{c.Parent, c.OnChildChange, transform, c.ParentPivot}}
}

func (c ParentAnchor) SetPivotPoint(pivot mgl32.Vec3) ParentAnchor {
	return ParentAnchor{&parentAnchor{c.Parent, c.OnChildChange, c.RelativeTransform, transform.NewPivotPoint(pivot)}}
}

package hierarchy

import (
	"engine/services/datastructures"
	"engine/services/ecs"
	"errors"
)

var (
	ErrParentCycle error = errors.New("parent cycle is not allowed")
)

type Component struct {
	Parent ecs.EntityID
}

func NewParent(parent ecs.EntityID) Component { return Component{parent} }

//

type ToolFactory ecs.ToolFactory[World, HierarchyTool]
type HierarchyTool interface {
	Hierarchy() Interface
}
type World interface {
	ecs.World
}
type Interface interface {
	Component() ecs.ComponentsArray[Component]

	// returns true if is child of any parent doesn't matter the depth
	IsChildOf(child ecs.EntityID, parent ecs.EntityID) bool
	SetParent(child ecs.EntityID, parent ecs.EntityID)
	Parent(child ecs.EntityID) (ecs.EntityID, bool)

	// from closest to furthest
	GetParents(child ecs.EntityID) datastructures.SparseSet[ecs.EntityID]
	GetOrderedParents(child ecs.EntityID) []ecs.EntityID

	// maintains order of children and adds component to children
	// even if children doesn't exist
	SetChildren(parent ecs.EntityID, children ...ecs.EntityID)

	Children(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
	// includes children of children
	FlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
}

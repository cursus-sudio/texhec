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

type Hierarchy interface {
	Hierarchy() Interface
}

type Interface interface {
	// returns true if is child of any parent doesn't matter the depth
	IsChildOf(child ecs.EntityID, parent ecs.EntityID) bool
	SetParent(child ecs.EntityID, parent ecs.EntityID)
	Parent(child ecs.EntityID) (ecs.EntityID, bool)

	// from closest to furthest
	GetParents(child ecs.EntityID) datastructures.SparseSet[ecs.EntityID]
	GetOrderedParents(child ecs.EntityID) []ecs.EntityID

	Children(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
	// includes children of children
	FlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
}

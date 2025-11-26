package hierarchy

import (
	"engine/services/datastructures"
	"engine/services/ecs"
)

type ParentComponent struct {
	Parent ecs.EntityID
}

func NewParent(parent ecs.EntityID) ParentComponent { return ParentComponent{parent} }

//

type Tool interface {
	Transaction() Transaction
	// parent can only be removed
	// OnParentRemove([]ecs.EntityID)
	// OnChildAppend([]ecs.EntityID)
	// OnChildRemove([]ecs.EntityID)
	// OnFlatChildAppend([]ecs.EntityID)
	// OnFlatChildRemove([]ecs.EntityID)
}

type Transaction interface {
	GetObject(ecs.EntityID) Object
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

// absolute components return errors only in case when
// entity is relative to parent that doesn't exist
type Object interface {
	Parent() ecs.EntityComponent[ParentComponent]

	// returns true if is child of any parent doesn't matter the depth
	IsChildOf(parent ecs.EntityID) bool
	// from closest to furthest
	GetParents() datastructures.SparseSet[ecs.EntityID]
	GetOrderedParents() []ecs.EntityID

	Children() datastructures.SparseSetReader[ecs.EntityID]
	// includes children of children
	FlatChildren() datastructures.SparseSetReader[ecs.EntityID]
}

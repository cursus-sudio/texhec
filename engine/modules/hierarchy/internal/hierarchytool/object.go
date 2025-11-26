package hierarchytool

import (
	"engine/modules/hierarchy"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type object struct {
	transaction

	parent ecs.EntityComponent[hierarchy.ParentComponent]
	entity ecs.EntityID
}

func newObject(
	t transaction,
	entity ecs.EntityID,
) hierarchy.Object {
	s := object{
		transaction: t,

		parent: t.parentTransaction.GetEntityComponent(entity),
		entity: entity,
	}
	return s
}

func (t object) Parent() ecs.EntityComponent[hierarchy.ParentComponent] { return t.parent }

func (t object) IsChildOf(wantedParent ecs.EntityID) bool {
	child := t.entity
	for {
		parent, err := t.transaction.GetObject(child).Parent().Get()
		if err != nil {
			return false
		}
		child = parent.Parent
		if child == wantedParent {
			return true
		}
	}
}

func (t object) GetParents() datastructures.SparseSet[ecs.EntityID] {
	orderedParents := t.GetOrderedParents()
	parents := datastructures.NewSparseSet[ecs.EntityID]()
	for _, parent := range orderedParents {
		parents.Add(parent)
	}
	return parents
}

func (t object) GetOrderedParents() []ecs.EntityID {
	parent, err := t.transaction.GetObject(t.entity).Parent().Get()
	if err != nil {
		return []ecs.EntityID{}
	}
	return t.tool.GetOrderedParents(parent)
}

func (t object) Children() datastructures.SparseSetReader[ecs.EntityID] {
	children, ok := t.parentChildren.Get(t.entity)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return children
}

func (t object) FlatChildren() datastructures.SparseSetReader[ecs.EntityID] {
	flatChildren, ok := t.tool.parentFlatChildren.Get(t.entity)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return flatChildren
}

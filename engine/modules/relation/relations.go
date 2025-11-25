package relation

import (
	"engine/services/datastructures"
	"engine/services/ecs"
)

type EntityToKeyTool[Key any] interface {
	Get(Key) (ecs.EntityID, bool)
	OnUpsert(func([]ecs.EntityID))
	OnRemove(func([]ecs.EntityID))
}

//

type ParentTool[Component any] interface {
	GetChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
	// notifies about parent change
	OnUpsert(newParentChildListener func([]ecs.EntityID))
	// notifies about parent removal
	OnRemove(parentListener func([]ecs.EntityID))
}

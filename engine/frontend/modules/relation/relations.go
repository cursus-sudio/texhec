package relation

import (
	"shared/services/datastructures"
	"shared/services/ecs"
)

type EntityToKeyTool[Key any] interface {
	Get(Key) (ecs.EntityID, bool)
	OnUpsert(func([]ecs.EntityID))
	OnRemove(func([]ecs.EntityID))
}

//

type EntityToEntitiesTool[Component any] interface {
	GetMany(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID]
	// notifies about parent change
	OnUpsert(newParentChildListener func([]ecs.EntityID))
	// notifies about parent removal
	OnRemove(parentListener func([]ecs.EntityID))
}

package indexing

import (
	"shared/services/datastructures"
	"shared/services/ecs"
)

type Indices[IndexType any] interface {
	Get(IndexType) (ecs.EntityID, bool)
	OnUpsert(func([]ecs.EntityID))
	OnRemove(func([]ecs.EntityID))
}

// type MapIndexTool[Component any] Indices[Component]

//

type ManyToOne[IndexType any] interface {
	GetMany(ecs.EntityID) datastructures.SparseSet[ecs.EntityID]
	GetOne(ecs.EntityID) (ecs.EntityID, bool)
	OnUpsert(func([]ecs.EntityID))
	OnRemove(func([]ecs.EntityID))
}

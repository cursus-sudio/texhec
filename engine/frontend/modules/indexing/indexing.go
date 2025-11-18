package indexing

import (
	"shared/services/ecs"
)

type Indices[IndexType any] interface {
	Get(IndexType) (ecs.EntityID, bool)
	OnUpsert(func([]ecs.EntityID))
	OnRemove(func([]ecs.EntityID))
}

type SpatialIndexTool[IndexType any] Indices[IndexType]
type MapIndexTool[Component any] Indices[Component]

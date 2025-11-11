package indexing

import (
	"shared/services/ecs"
)

type Indices[IndexType any] interface {
	Get(IndexType) (ecs.EntityID, bool)
}

type SpatialIndexTool[Component, IndexType any] Indices[IndexType]
type MapIndexTool[Component any] Indices[Component]

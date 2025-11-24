package onetokey

import (
	"frontend/modules/relation"
	"shared/services/ecs"
)

func NewSpatialRelationFactory[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndexFactory func(ecs.World) func(ecs.EntityID) (IndexType, bool),
	indexNumber func(IndexType) uint32,
) ecs.ToolFactory[relation.EntityToKeyTool[IndexType]] {
	return ecs.NewToolFactory(func(w ecs.World) relation.EntityToKeyTool[IndexType] {
		if index, err := ecs.GetGlobal[spatialRelation[IndexType]](w); err == nil {
			return index
		}
		w.LockGlobals()
		defer w.UnlockGlobals()
		if index, err := ecs.GetGlobal[spatialRelation[IndexType]](w); err == nil {
			return index
		}
		query := queryFactory(w)
		componentIndex := componentIndexFactory(w)
		return newIndex(w, query, componentIndex, indexNumber)
	})
}

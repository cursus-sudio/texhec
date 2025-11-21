package manytoone

import (
	"frontend/modules/relation"
	"shared/services/ecs"
)

func NewManyToOneFactory[Component any](componentParent func(Component) ecs.EntityID) ecs.ToolFactory[relation.EntityToEntitiesTool[Component]] {
	return ecs.NewToolFactory(func(w ecs.World) relation.EntityToEntitiesTool[Component] {
		if index, err := ecs.GetGlobal[manyToOne[Component]](w); err == nil {
			return index
		}
		w.LockGlobals()
		defer w.UnlockGlobals()
		if index, err := ecs.GetGlobal[manyToOne[Component]](w); err == nil {
			return index
		}
		return newIndex(w, componentParent)
	})
}

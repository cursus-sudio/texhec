package parent

import (
	"frontend/modules/relation"
	"shared/services/ecs"
)

func NewParentToolFactory[Component any](componentParent func(Component) ecs.EntityID) ecs.ToolFactory[relation.ParentTool[Component]] {
	return ecs.NewToolFactory(func(w ecs.World) relation.ParentTool[Component] {
		if index, err := ecs.GetGlobal[parent[Component]](w); err == nil {
			return index
		}
		w.LockGlobals()
		defer w.UnlockGlobals()
		if index, err := ecs.GetGlobal[parent[Component]](w); err == nil {
			return index
		}
		return newIndex(w, componentParent)
	})
}

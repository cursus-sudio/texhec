package internal

import (
	"engine/modules/relation"
	"engine/modules/uuid"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	relation.EntityToKeyTool[uuid.UUID]
}

func NewToolFactory(
	toolFactory ecs.ToolFactory[relation.EntityToKeyTool[uuid.UUID]],
) ecs.ToolFactory[uuid.Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) uuid.Tool {
		if t, err := ecs.GetGlobal[tool](w); err == nil {
			return t
		}
		mutex.Lock()
		defer mutex.Unlock()
		if t, err := ecs.GetGlobal[tool](w); err == nil {
			return t
		}
		t := tool{
			toolFactory.Build(w),
		}
		w.SaveGlobal(t)
		return t

	})
}

func (t tool) Entity(uuid uuid.UUID) (ecs.EntityID, bool) {
	return t.Get(uuid)
}

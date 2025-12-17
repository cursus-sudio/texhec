package internal

import (
	"engine/modules/relation"
	"engine/modules/uuid"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	relation.EntityToKeyTool[uuid.UUID]
	uuid.Factory
}

func NewToolFactory(
	toolFactory ecs.ToolFactory[relation.EntityToKeyTool[uuid.UUID]],
	uuidFactory uuid.Factory,
) ecs.ToolFactory[uuid.UUIDTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) uuid.UUIDTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			toolFactory.Build(w),
			uuidFactory,
		}
		w.SaveGlobal(t)
		return t

	})
}

func (t tool) Entity(uuid uuid.UUID) (ecs.EntityID, bool) {
	return t.Get(uuid)
}

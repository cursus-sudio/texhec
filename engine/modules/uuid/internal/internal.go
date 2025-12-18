package internal

import (
	"engine/modules/relation"
	"engine/modules/uuid"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	uuidArray ecs.ComponentsArray[uuid.Component]
	relation.EntityToKeyTool[uuid.UUID]
	uuid.Factory
}

func NewToolFactory(
	toolFactory ecs.ToolFactory[ecs.World, relation.EntityToKeyTool[uuid.UUID]],
	uuidFactory uuid.Factory,
) ecs.ToolFactory[uuid.World, uuid.UUIDTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w uuid.World) uuid.UUIDTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			ecs.GetComponentsArray[uuid.Component](w),
			toolFactory.Build(w),
			uuidFactory,
		}
		w.SaveGlobal(t)
		return t

	})
}

func (t tool) UUID() uuid.Interface { return t }

func (t tool) Component() ecs.ComponentsArray[uuid.Component] { return t.uuidArray }

func (t tool) Entity(uuid uuid.UUID) (ecs.EntityID, bool) {
	return t.Get(uuid)
}

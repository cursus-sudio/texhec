package internal

import (
	"engine/modules/groups"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

type tool struct {
	logger logger.Logger

	world groups.World

	inheritArray ecs.ComponentsArray[groups.InheritGroupsComponent]
	groupsArray  ecs.ComponentsArray[groups.GroupsComponent]
}

func NewToolFactory(logger logger.Logger) groups.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w groups.World) groups.GroupsTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := &tool{
			logger: logger,

			world:        w,
			inheritArray: ecs.GetComponentsArray[groups.InheritGroupsComponent](w),
			groupsArray:  ecs.GetComponentsArray[groups.GroupsComponent](w),
		}
		w.SaveGlobal(t)
		t.Init()
		return t
	})
}

func (t *tool) Groups() groups.Interface {
	return t
}

func (t *tool) Component() ecs.ComponentsArray[groups.GroupsComponent] {
	return t.groupsArray
}
func (t *tool) Inherit() ecs.ComponentsArray[groups.InheritGroupsComponent] {
	return t.inheritArray
}

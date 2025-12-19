package tool

import (
	"engine/modules/transition"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	easing ecs.ComponentsArray[transition.EasingComponent]
}

func NewToolFactory() transition.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w transition.World) transition.TransitionTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			ecs.GetComponentsArray[transition.EasingComponent](w),
		}

		w.SaveGlobal(t)
		return t
	})
}

func (t tool) Transition() transition.Interface {
	return t
}

func (t tool) Easing() ecs.ComponentsArray[transition.EasingComponent] {
	return t.easing
}

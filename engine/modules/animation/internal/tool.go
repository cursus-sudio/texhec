package internal

import (
	"engine/modules/animation"
	"engine/services/ecs"
	"sync"
)

// type AnimationTool interface {
// 	Animation() Interface
// }
// type World interface {
// 	ecs.World
// }
// type Interface interface {
// 	Component() ecs.ComponentsArray[AnimationComponent]
// }

type tool struct {
	world           animation.World
	animationArrary ecs.ComponentsArray[animation.AnimationComponent]
	loopArrary      ecs.ComponentsArray[animation.LoopComponent]
}

func NewToolFactory() animation.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w animation.World) animation.AnimationTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			w,
			ecs.GetComponentsArray[animation.AnimationComponent](w),
			ecs.GetComponentsArray[animation.LoopComponent](w),
		}
		w.SaveGlobal(t)
		return t
	})
}

func (t tool) Animation() animation.Interface {
	return t
}
func (t tool) Component() ecs.ComponentsArray[animation.AnimationComponent] {
	return t.animationArrary
}
func (t tool) Loop() ecs.ComponentsArray[animation.LoopComponent] {
	return t.loopArrary
}

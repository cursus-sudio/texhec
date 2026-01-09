package renderer

import (
	"engine/modules/genericrenderer"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	genericrenderer.World

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewToolFactory() genericrenderer.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w genericrenderer.World) genericrenderer.GenericRendererTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := &tool{
			w,
			ecs.GetComponentsArray[genericrenderer.PipelineComponent](w),
		}
		t.SaveGlobal(t)
		return t
	})
}

func (t *tool) GenericRenderer() genericrenderer.Interface {
	return t
}

func (t *tool) Pipeline() ecs.ComponentsArray[genericrenderer.PipelineComponent] {
	return t.pipelineArray
}

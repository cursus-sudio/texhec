package renderer

import (
	"engine/modules/genericrenderer"
	"engine/services/ecs"
)

type tool struct {
	World ecs.World

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewService(
	world ecs.World,
) genericrenderer.Service {
	t := &tool{
		world,
		ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
	}
	return t
}

func (t *tool) Pipeline() ecs.ComponentsArray[genericrenderer.PipelineComponent] {
	return t.pipelineArray
}

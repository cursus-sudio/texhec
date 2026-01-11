package renderer

import (
	"engine/modules/genericrenderer"
	"engine/services/ecs"
)

type service struct {
	World ecs.World

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewService(
	world ecs.World,
) genericrenderer.Service {
	t := &service{
		world,
		ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
	}
	return t
}

func (t *service) Pipeline() ecs.ComponentsArray[genericrenderer.PipelineComponent] {
	return t.pipelineArray
}

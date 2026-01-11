package renderer

import (
	"engine/modules/genericrenderer"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World ecs.World `inject:"1"`

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewService(c ioc.Dic) genericrenderer.Service {
	t := ioc.GetServices[*service](c)
	t.pipelineArray = ecs.GetComponentsArray[genericrenderer.PipelineComponent](t.World)
	return t
}

func (t *service) Pipeline() ecs.ComponentsArray[genericrenderer.PipelineComponent] {
	return t.pipelineArray
}

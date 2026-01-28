package direct

import (
	"engine/modules/genericrenderer"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World ecs.World `inject:"1"`

	pipelineArray ecs.ComponentsArray[genericrenderer.DirectComponent]
}

func NewService(c ioc.Dic) genericrenderer.Service {
	t := ioc.GetServices[*service](c)
	t.pipelineArray = ecs.GetComponentsArray[genericrenderer.DirectComponent](t.World)
	return t
}

func (t *service) Direct() ecs.ComponentsArray[genericrenderer.DirectComponent] {
	return t.pipelineArray
}

package service

import (
	"engine/modules/renderer"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World ecs.World `inject:"1"`

	pipelineArray ecs.ComponentsArray[renderer.DirectComponent]
}

func NewService(c ioc.Dic) renderer.Service {
	t := ioc.GetServices[*service](c)
	t.pipelineArray = ecs.GetComponentsArray[renderer.DirectComponent](t.World)
	return t
}

func (t *service) Render(entity ecs.EntityID) {
	t.Direct().Set(entity, renderer.DirectComponent{})
}

func (t *service) Direct() ecs.ComponentsArray[renderer.DirectComponent] {
	return t.pipelineArray
}

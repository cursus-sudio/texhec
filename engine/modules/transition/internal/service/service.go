package service

import (
	"engine/modules/transition"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World  ecs.World `inject:"1"`
	easing ecs.ComponentsArray[transition.EasingComponent]
}

func NewService(c ioc.Dic) transition.Service {
	t := ioc.GetServices[*service](c)
	t.easing = ecs.GetComponentsArray[transition.EasingComponent](t.World)

	return t
}

func (t *service) Easing() ecs.ComponentsArray[transition.EasingComponent] {
	return t.easing
}

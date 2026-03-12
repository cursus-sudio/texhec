package service

import (
	"engine/modules/transition"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World          ecs.World `inject:"1"`
	easing         ecs.ComponentsArray[transition.EasingComponent]
	easingFunction ecs.ComponentsArray[transition.EasingFunctionComponent]
}

func NewService(c ioc.Dic) transition.Service {
	t := ioc.GetServices[*service](c)
	t.easing = ecs.GetComponentsArray[transition.EasingComponent](t.World)
	t.easingFunction = ecs.GetComponentsArray[transition.EasingFunctionComponent](t.World)

	return t
}

func (t *service) Easing() ecs.ComponentsArray[transition.EasingComponent] {
	return t.easing
}
func (t *service) EasingFunction() ecs.ComponentsArray[transition.EasingFunctionComponent] {
	return t.easingFunction
}

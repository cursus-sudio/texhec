package service

import (
	"engine/modules/transition"
	"engine/services/ecs"
)

type service struct {
	easing ecs.ComponentsArray[transition.EasingComponent]
}

func NewService(w ecs.World) transition.Service {
	t := &service{
		ecs.GetComponentsArray[transition.EasingComponent](w),
	}

	return t
}

func (t *service) Easing() ecs.ComponentsArray[transition.EasingComponent] {
	return t.easing
}

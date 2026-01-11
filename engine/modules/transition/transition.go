package transition

import (
	"engine/services/ecs"
	"time"
)

type System ecs.SystemRegister

type Service interface {
	Easing() ecs.ComponentsArray[EasingComponent]
}

//

type Lerp[Component any] interface {
	Lerp(Component, float32) Component
}

//

type Progress float32

type TransitionComponent[Component Lerp[Component]] struct {
	From, To Component
	Progress,
	Duration time.Duration
}

func NewTransition[Component Lerp[Component]](
	from, to Component,
	duration time.Duration,
) TransitionComponent[Component] {
	return TransitionComponent[Component]{
		From:     from,
		To:       to,
		Progress: 0,
		Duration: duration,
	}
}

//

// saves transition component
type TransitionEvent[Component Lerp[Component]] struct {
	Entity    ecs.EntityID
	Component TransitionComponent[Component]
}

func NewTransitionEvent[Component Lerp[Component]](
	entity ecs.EntityID,
	from, to Component,
	duration time.Duration,
) TransitionEvent[Component] {
	return TransitionEvent[Component]{
		Entity: entity,
		Component: NewTransition(
			from, to,
			duration,
		),
	}
}

//

type EasingID uint16
type EasingFunction func(t Progress) Progress

type EasingService interface {
	Set(EasingID, EasingFunction)
	Get(EasingID) (EasingFunction, bool)
}

type EasingComponent struct {
	ID EasingID
}

func NewEasing(id EasingID) EasingComponent {
	return EasingComponent{id}
}

package transition

import (
	"engine/services/ecs"
	"time"
)

type System ecs.SystemRegister[World]

type ToolFactory ecs.ToolFactory[World, TransitionTool]
type TransitionTool interface {
	Transition() Interface
}
type World interface {
	ecs.World
}
type Interface interface {
	Easing() ecs.ComponentsArray[EasingComponent]
}

//

type BlendableComponent[Component any] interface {
	Blend(Component, float32) Component
}

//

type Progress float32

type TransitionComponent[Component BlendableComponent[Component]] struct {
	From, To Component
	Progress,
	Duration time.Duration
}

func NewTransition[Component BlendableComponent[Component]](
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
type TransitionEvent[Component BlendableComponent[Component]] struct {
	Entity    ecs.EntityID
	Component TransitionComponent[Component]
}

func NewTransitionEvent[Component BlendableComponent[Component]](
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

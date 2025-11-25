package internal

import (
	"engine/modules/animation"
	"engine/services/ecs"
)

type Transition struct {
	From, To           any // this should be of type T
	Start, End         animation.AnimationState
	EasingFunction     animation.EasingFunction
	TransitionFunction animation.AnyTransitionFunction
}

func (t Transition) Duration() animation.AnimationState {
	return t.End - t.Start
}

func (t Transition) NormalizedState(state animation.AnimationState) animation.AnimationState {
	duration := t.Duration()
	state = state - t.Start
	if state >= duration {
		return 1
	} else if state <= 0 {
		return 0
	}
	return state / duration
}

func (t Transition) CallTransitionFunction(entity ecs.EntityID, normalizedState animation.AnimationState) error {
	args := animation.TransitionFunctionAnyArgument{
		Entity: entity,
		From:   t.From,
		To:     t.To,
		State:  normalizedState,
	}
	return t.TransitionFunction(args)
}

type Animation struct {
	Events      []animation.Event
	Transitions []Transition
}

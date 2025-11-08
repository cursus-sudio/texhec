package internal

import (
	"fmt"
	"frontend/modules/animation"
	"frontend/services/frames"
	"reflect"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/events"
)

type system struct {
	easingFunctions datastructures.SparseArray[animation.EasingFunctionID, animation.EasingFunction]
	animations      datastructures.SparseArray[animation.AnimationID, Animation]

	logger                logger.Logger
	world                 ecs.World
	animationsArray       ecs.ComponentsArray[animation.AnimationComponent]
	animationsTransaction ecs.ComponentsArrayTransaction[animation.AnimationComponent]
}

func NewSystem(
	b AnimationSystemBuilder,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transitionFunctions := make(map[reflect.Type]animation.AnyTransitionFunction, len(b.transitionFunctions))
		for key, transitionFunction := range b.transitionFunctions {
			transitionFunctions[key] = transitionFunction(w)
		}
		animations := datastructures.NewSparseArray[animation.AnimationID, Animation]()
		for _, id := range b.animations.GetIndices() {
			animationData, _ := b.animations.Get(id)
			animation := Animation{
				Events:      animationData.Events,
				Transitions: make([]Transition, 0, len(animationData.Transitions)),
			}

			for _, transition := range animationData.Transitions {
				easingFunction, ok := b.easingFunctions.Get(transition.EasingFunction)
				if !ok {
					b.logger.Error(fmt.Errorf(
						"expected easing function with id \"%v\" to exist. skipping transition",
						transition.EasingFunction,
					))
					continue
				}
				animation.Transitions = append(animation.Transitions, Transition{
					From:  transition.From,
					To:    transition.To,
					Start: transition.Start,
					End:   transition.End,

					EasingFunction:     easingFunction,
					TransitionFunction: transitionFunctions[reflect.TypeOf(transition.From)],
				})
			}

			animations.Set(id, animation)
		}

		animationsArray := ecs.GetComponentsArray[animation.AnimationComponent](w.Components())
		s := &system{
			easingFunctions:       b.easingFunctions,
			animations:            animations,
			world:                 w,
			animationsArray:       animationsArray,
			animationsTransaction: animationsArray.Transaction(),
		}

		events.ListenE(w.EventsBuilder(), s.ListenE)

		return nil
	})
}

func (s *system) ListenE(event frames.FrameEvent) error {
	for _, entity := range s.animationsArray.GetEntities() {
		animationComp, err := s.animationsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		animationComp.AddElapsedTime(event.Delta)
		s.animationsTransaction.SaveComponent(entity, animationComp)

		if animationComp.PreviousState == animationComp.State {
			continue
		}

		animationData, ok := s.animations.Get(animationComp.AnimationID)
		if !ok {
			s.logger.Error(fmt.Errorf(
				"expected animation data with id \"%v\" to exist",
				animationComp.AnimationID,
			))
			continue
		}

		for _, eventData := range animationData.Events {
			if animationComp.PreviousState < eventData.Trigger &&
				eventData.Trigger < animationComp.State {
				events.EmitAny(s.world.Events(), eventData.Event)
			}
		}

		for _, transition := range animationData.Transitions {
			previous := transition.NormalizedState(animationComp.PreviousState)
			currentState := transition.NormalizedState(animationComp.State)
			if previous == currentState {
				continue
			}
			currentState = transition.EasingFunction(currentState)
			transition.CallTransitionFunction(entity, currentState)
		}
	}
	return s.animationsTransaction.Flush()
}

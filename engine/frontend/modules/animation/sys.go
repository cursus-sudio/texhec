package animation

import (
	"reflect"
	"shared/services/ecs"
)

type System ecs.SystemRegister

//

type TransitionFunctionArgument[Component any] struct {
	Entity   ecs.EntityID
	From, To Component
	State    AnimationState
}

type TransitionFunctionAnyArgument struct {
	Entity   ecs.EntityID
	From, To any // this should be of type T
	State    AnimationState
}

type TransitionFunction[Component any] func(arg TransitionFunctionArgument[Component]) error
type AnyTransitionFunction func(arg TransitionFunctionAnyArgument) error

func AddTransitionFunction[Component any](
	b AnimationSystemBuilder,
	transitionFunction func(ecs.World) TransitionFunction[Component],
) {
	b.AddTransitionFunction(reflect.TypeFor[Component](), func(w ecs.World) AnyTransitionFunction {
		inner := transitionFunction(w)
		return func(anyArg TransitionFunctionAnyArgument) error {
			arg := TransitionFunctionArgument[Component]{
				Entity: anyArg.Entity,
				From:   anyArg.From.(Component),
				To:     anyArg.To.(Component),
				State:  anyArg.State,
			}
			return inner(arg)
		}
	})
}

// Build method is private and used internally.
// If 2 times you'll save something with the same ID only newest will be considered
type AnimationSystemBuilder interface {
	AddEasingFunction(EasingFunctionID, EasingFunction)
	AddTransitionFunction(t reflect.Type, transitionFunction func(ecs.World) AnyTransitionFunction)
	AddAnimation(AnimationID, Animation)
}

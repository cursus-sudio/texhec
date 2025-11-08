package internal

import (
	"frontend/modules/animation"
	"reflect"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"
)

type AnimationSystemBuilder struct {
	logger              logger.Logger
	easingFunctions     datastructures.SparseArray[animation.EasingFunctionID, animation.EasingFunction]
	transitionFunctions map[reflect.Type]func(ecs.World) animation.AnyTransitionFunction
	animations          datastructures.SparseArray[animation.AnimationID, animation.Animation]
}

func NewBuilder() AnimationSystemBuilder {
	return AnimationSystemBuilder{
		easingFunctions:     datastructures.NewSparseArray[animation.EasingFunctionID, animation.EasingFunction](),
		transitionFunctions: make(map[reflect.Type]func(ecs.World) animation.AnyTransitionFunction),
		animations:          datastructures.NewSparseArray[animation.AnimationID, animation.Animation](),
	}
}

func (b AnimationSystemBuilder) AddEasingFunction(
	easingFunctionID animation.EasingFunctionID,
	easingFunction animation.EasingFunction,
) {
	b.easingFunctions.Set(easingFunctionID, easingFunction)
}

func (b AnimationSystemBuilder) AddTransitionFunction(
	t reflect.Type,
	transitionFunction func(ecs.World) animation.AnyTransitionFunction,
) {
	b.transitionFunctions[t] = transitionFunction
}

func (b AnimationSystemBuilder) AddAnimation(
	animationID animation.AnimationID,
	animationData animation.Animation,
) {
	b.animations.Set(animationID, animationData)
}

func (b AnimationSystemBuilder) Build() ecs.SystemRegister {
	return NewSystem(b)
}

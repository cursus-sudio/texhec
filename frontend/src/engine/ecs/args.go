package ecs

import . "frontend/src/engine/ecs/ecsargs"

type Args struct {
	DeltaTime DeltaTime
}

func NewArgs(
	deltaTime DeltaTime,
) Args {
	return Args{
		DeltaTime: deltaTime}
}

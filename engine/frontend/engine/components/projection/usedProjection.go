package projection

import (
	"errors"
	"shared/services/ecs"
)

var (
	ErrExpectedUsedProjectionToImplementProjection error = errors.New("expected component type which implements Projection interface")
)

type UsedProjection struct {
	ProjectionComponent ecs.ComponentType
}

func NewUsedProjection[P any]() UsedProjection {
	var pZero P
	return UsedProjection{ProjectionComponent: ecs.GetComponentType(pZero)}
}

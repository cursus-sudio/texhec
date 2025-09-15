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

func NewUsedProjection[P Projection]() UsedProjection {
	var pZero P
	return UsedProjection{ProjectionComponent: ecs.GetComponentType(pZero)}
}

func (usedProjection UsedProjection) GetCameraProjection(world ecs.World, camera ecs.EntityID) (Projection, error) {
	anyProj, err := world.GetAnyComponent(camera, usedProjection.ProjectionComponent)
	if err != nil {
		return nil, err
	}

	proj, ok := anyProj.(Projection)
	if !ok {
		return nil, ErrExpectedUsedProjectionToImplementProjection
	}

	return proj, nil
}

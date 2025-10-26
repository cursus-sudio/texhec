package cameras

import (
	"errors"
	"fmt"
	cameracomponent "frontend/engine/components/camera"
	"shared/services/ecs"
)

var ErrMissingConstructor error = errors.New("missing camera type constructor")

type CameraConstructors interface {
	Get(ecs.World, ecs.EntityID) (Camera, error)
}

type cameraConstructors struct {
	constructors map[ecs.ComponentType]func(ecs.World, ecs.EntityID) (Camera, error)
}

func (c *cameraConstructors) Get(world ecs.World, entity ecs.EntityID) (Camera, error) {
	cameraArray := ecs.GetComponentsArray[cameracomponent.Camera](world.Components())
	cameraComponent, err := cameraArray.GetComponent(entity)
	if err != nil {
		return nil, err
	}
	constructor, ok := c.constructors[cameraComponent.Projection]
	if !ok {
		return nil, errors.Join(
			ErrMissingConstructor,
			fmt.Errorf("missing constructor for \"%s\"", cameraComponent.Projection.String()),
		)
	}

	camera, err := constructor(world, entity)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

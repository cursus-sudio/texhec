package cameras

import (
	"errors"
	"fmt"
	cameracomponent "frontend/engine/components/camera"
	"shared/services/ecs"
)

var ErrMissingConstructor error = errors.New("missing camera type constructor")

type CameraResolver interface {
	Get(ecs.EntityID) (Camera, error)
}

type cameraResolver struct {
	cameraArray  ecs.ComponentsArray[cameracomponent.Camera]
	constructors map[ecs.ComponentType]func(ecs.EntityID) (Camera, error)
}

func (c *cameraResolver) Get(entity ecs.EntityID) (Camera, error) {
	cameraComponent, err := c.cameraArray.GetComponent(entity)
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

	camera, err := constructor(entity)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

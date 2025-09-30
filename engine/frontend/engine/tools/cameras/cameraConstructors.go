package cameras

import (
	"errors"
	"fmt"
	"shared/services/ecs"
)

var ErrMissingConstructor error = errors.New("missing camera type constructor")

type CameraConstructors interface {
	Get(ecs.EntityID, ecs.ComponentType) (Camera, error)
}

type cameraConstructors struct {
	constructors map[ecs.ComponentType]func(ecs.EntityID) (Camera, error)
}

func (c *cameraConstructors) Get(entity ecs.EntityID, componentType ecs.ComponentType) (Camera, error) {
	constructor, ok := c.constructors[componentType]
	if !ok {
		return nil, errors.Join(
			ErrMissingConstructor,
			fmt.Errorf("missing constructor for \"%s\"", componentType.String()),
		)
	}

	camera, err := constructor(entity)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

package cameratool

import (
	"errors"
	"fmt"
	"frontend/modules/camera"
	"shared/services/ecs"
)

type cameraResolver struct {
	cameraArray  ecs.ComponentsArray[camera.CameraComponent]
	constructors map[ecs.ComponentType]func(ecs.EntityID) (camera.CameraService, error)
}

func (c *cameraResolver) Get(entity ecs.EntityID) (camera.CameraService, error) {
	cameraComponent, err := c.cameraArray.GetComponent(entity)
	if err != nil {
		return nil, errors.Join(
			camera.ErrNotCamera,
			err,
		)
	}
	constructor, ok := c.constructors[cameraComponent.Projection]
	if !ok {
		return nil, errors.Join(
			camera.ErrNotCamera,
			fmt.Errorf("missing constructor for \"%s\"", cameraComponent.Projection.String()),
		)
	}

	camera, err := constructor(entity)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

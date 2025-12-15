package cameratool

import (
	"engine/modules/camera"
	"engine/services/ecs"
	"errors"
	"fmt"
	"reflect"
)

type cameraResolver struct {
	cameraArray  ecs.ComponentsArray[camera.CameraComponent]
	constructors map[reflect.Type]func(ecs.EntityID) (camera.Object, error)
}

func (c *cameraResolver) Camera() camera.Interface { return c }

func (c *cameraResolver) GetObject(entity ecs.EntityID) (camera.Object, error) {
	cameraComponent, ok := c.cameraArray.GetComponent(entity)
	if !ok {
		return nil, camera.ErrNotCamera
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

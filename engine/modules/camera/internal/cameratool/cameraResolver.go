package cameratool

import (
	"engine/modules/camera"
	"engine/services/ecs"
	"errors"
	"fmt"
	"reflect"
)

type cameraResolver struct {
	cameraArray  ecs.ComponentsArray[camera.Component]
	constructors map[reflect.Type]func(ecs.EntityID) (camera.Object, error)

	mobileCamera       ecs.ComponentsArray[camera.MobileCameraComponent]
	cameraLimits       ecs.ComponentsArray[camera.CameraLimitsComponent]
	viewport           ecs.ComponentsArray[camera.ViewportComponent]
	normalizedViewport ecs.ComponentsArray[camera.NormalizedViewportComponent]

	ortho              ecs.ComponentsArray[camera.OrthoComponent]
	orthoResolution    ecs.ComponentsArray[camera.OrthoResolutionComponent]
	perspective        ecs.ComponentsArray[camera.PerspectiveComponent]
	dynamicPerspective ecs.ComponentsArray[camera.DynamicPerspectiveComponent]
}

func (c *cameraResolver) Camera() camera.Interface { return c }

func (c *cameraResolver) Component() ecs.ComponentsArray[camera.Component] {
	return c.cameraArray
}

func (c *cameraResolver) GetObject(entity ecs.EntityID) (camera.Object, error) {
	cameraComponent, ok := c.cameraArray.Get(entity)
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

func (c *cameraResolver) Mobile() ecs.ComponentsArray[camera.MobileCameraComponent] {
	return c.mobileCamera
}
func (c *cameraResolver) Limits() ecs.ComponentsArray[camera.CameraLimitsComponent] {
	return c.cameraLimits
}
func (c *cameraResolver) Viewport() ecs.ComponentsArray[camera.ViewportComponent] {
	return c.viewport
}
func (c *cameraResolver) NormalizedViewport() ecs.ComponentsArray[camera.NormalizedViewportComponent] {
	return c.normalizedViewport
}

func (c *cameraResolver) Ortho() ecs.ComponentsArray[camera.OrthoComponent] {
	return c.ortho
}
func (c *cameraResolver) OrthoResolution() ecs.ComponentsArray[camera.OrthoResolutionComponent] {
	return c.orthoResolution
}
func (c *cameraResolver) Perspective() ecs.ComponentsArray[camera.PerspectiveComponent] {
	return c.perspective
}
func (c *cameraResolver) DynamicPerspective() ecs.ComponentsArray[camera.DynamicPerspectiveComponent] {
	return c.dynamicPerspective
}

package cameratool

import (
	"engine/modules/camera"
	"engine/services/ecs"
	"reflect"
)

type CameraResolverFactory interface {
	Register(
		reflect.Type,
		func(camera.World) func(ecs.EntityID) (camera.Object, error),
	)
	ecs.ToolFactory[camera.World, camera.CameraTool]
}

type cameraResolverFactory struct {
	constructors map[reflect.Type]func(camera.World) func(ecs.EntityID) (camera.Object, error)
}

func NewCameraResolverFactory() CameraResolverFactory {
	return &cameraResolverFactory{
		constructors: make(map[reflect.Type]func(camera.World) func(ecs.EntityID) (camera.Object, error)),
	}
}

func (f *cameraResolverFactory) Register(
	componentType reflect.Type,
	ctor func(camera.World) func(ecs.EntityID) (camera.Object, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraResolverFactory) Build(world camera.World) camera.CameraTool {
	ctors := make(map[reflect.Type]func(ecs.EntityID) (camera.Object, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraResolver{
		cameraArray:  ecs.GetComponentsArray[camera.Component](world),
		constructors: ctors,

		mobileCamera:       ecs.GetComponentsArray[camera.MobileCameraComponent](world),
		cameraLimits:       ecs.GetComponentsArray[camera.CameraLimitsComponent](world),
		viewport:           ecs.GetComponentsArray[camera.ViewportComponent](world),
		normalizedViewport: ecs.GetComponentsArray[camera.NormalizedViewportComponent](world),

		ortho:              ecs.GetComponentsArray[camera.OrthoComponent](world),
		orthoResolution:    ecs.GetComponentsArray[camera.OrthoResolutionComponent](world),
		perspective:        ecs.GetComponentsArray[camera.PerspectiveComponent](world),
		dynamicPerspective: ecs.GetComponentsArray[camera.DynamicPerspectiveComponent](world),
	}
}

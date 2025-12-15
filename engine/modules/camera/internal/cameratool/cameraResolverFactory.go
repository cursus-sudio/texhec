package cameratool

import (
	"engine/modules/camera"
	"engine/services/ecs"
	"reflect"
)

type CameraResolverFactory interface {
	Register(
		reflect.Type,
		func(ecs.World) func(ecs.EntityID) (camera.Object, error),
	)
	ecs.ToolFactory[camera.Camera]
}

type cameraResolverFactory struct {
	constructors map[reflect.Type]func(ecs.World) func(ecs.EntityID) (camera.Object, error)
}

func NewCameraResolverFactory() CameraResolverFactory {
	return &cameraResolverFactory{
		constructors: make(map[reflect.Type]func(ecs.World) func(ecs.EntityID) (camera.Object, error)),
	}
}

func (f *cameraResolverFactory) Register(
	componentType reflect.Type,
	ctor func(ecs.World) func(ecs.EntityID) (camera.Object, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraResolverFactory) Build(world ecs.World) camera.Camera {
	ctors := make(map[reflect.Type]func(ecs.EntityID) (camera.Object, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraResolver{
		cameraArray:  ecs.GetComponentsArray[camera.CameraComponent](world),
		constructors: ctors,
	}
}

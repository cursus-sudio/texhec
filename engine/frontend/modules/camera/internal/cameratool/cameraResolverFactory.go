package cameratool

import (
	"frontend/modules/camera"
	"shared/services/ecs"
)

type CameraResolverFactory interface {
	Register(
		ecs.ComponentType,
		func(ecs.World) func(ecs.EntityID) (camera.CameraService, error),
	)
	ecs.ToolFactory[camera.CameraTool]
}

type cameraResolverFactory struct {
	constructors map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (camera.CameraService, error)
}

func NewCameraResolverFactory() CameraResolverFactory {
	return &cameraResolverFactory{
		constructors: make(map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (camera.CameraService, error)),
	}
}

func (f *cameraResolverFactory) Register(
	componentType ecs.ComponentType,
	ctor func(ecs.World) func(ecs.EntityID) (camera.CameraService, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraResolverFactory) Build(world ecs.World) camera.CameraTool {
	ctors := make(map[ecs.ComponentType]func(ecs.EntityID) (camera.CameraService, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraResolver{
		cameraArray:  ecs.GetComponentsArray[camera.CameraComponent](world),
		constructors: ctors,
	}
}

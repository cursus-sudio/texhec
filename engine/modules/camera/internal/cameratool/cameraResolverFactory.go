package cameratool

import (
	"engine/modules/camera"
	"engine/services/ecs"
)

type CameraResolverFactory interface {
	Register(
		ecs.ComponentType,
		func(ecs.World) func(ecs.EntityID) (camera.Object, error),
	)
	ecs.ToolFactory[camera.Tool]
}

type cameraResolverFactory struct {
	constructors map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (camera.Object, error)
}

func NewCameraResolverFactory() CameraResolverFactory {
	return &cameraResolverFactory{
		constructors: make(map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (camera.Object, error)),
	}
}

func (f *cameraResolverFactory) Register(
	componentType ecs.ComponentType,
	ctor func(ecs.World) func(ecs.EntityID) (camera.Object, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraResolverFactory) Build(world ecs.World) camera.Tool {
	ctors := make(map[ecs.ComponentType]func(ecs.EntityID) (camera.Object, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraResolver{
		cameraArray:  ecs.GetComponentsArray[camera.CameraComponent](world),
		constructors: ctors,
	}
}

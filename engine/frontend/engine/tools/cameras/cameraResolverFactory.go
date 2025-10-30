package cameras

import (
	cameracomponent "frontend/engine/components/camera"
	"shared/services/ecs"
)

type CameraResolverFactory interface {
	Register(
		ecs.ComponentType,
		func(ecs.World) func(ecs.EntityID) (Camera, error),
	)
	ecs.ToolFactory[CameraResolver]
}

type cameraResolverFactory struct {
	constructors map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (Camera, error)
}

func (f *cameraResolverFactory) Register(
	componentType ecs.ComponentType,
	ctor func(ecs.World) func(ecs.EntityID) (Camera, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraResolverFactory) Build(world ecs.World) CameraResolver {
	ctors := make(map[ecs.ComponentType]func(ecs.EntityID) (Camera, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraResolver{
		cameraArray:  ecs.GetComponentsArray[cameracomponent.Camera](world.Components()),
		constructors: ctors,
	}
}

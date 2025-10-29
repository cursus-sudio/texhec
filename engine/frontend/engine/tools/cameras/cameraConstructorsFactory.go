package cameras

import (
	cameracomponent "frontend/engine/components/camera"
	"shared/services/ecs"
)

type CameraConstructorsFactory interface {
	Register(
		ecs.ComponentType,
		func(ecs.World) func(ecs.EntityID) (Camera, error),
	)
	ecs.ToolFactory[CameraConstructors]
}

type cameraConstructorsFactory struct {
	constructors map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (Camera, error)
}

func (f *cameraConstructorsFactory) Register(
	componentType ecs.ComponentType,
	ctor func(ecs.World) func(ecs.EntityID) (Camera, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraConstructorsFactory) Build(world ecs.World) CameraConstructors {
	ctors := make(map[ecs.ComponentType]func(ecs.EntityID) (Camera, error))
	for key, ctor := range f.constructors {
		ctors[key] = ctor(world)
	}
	return &cameraConstructors{
		cameraArray:  ecs.GetComponentsArray[cameracomponent.Camera](world.Components()),
		constructors: ctors,
	}
}

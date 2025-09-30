package cameras

import "shared/services/ecs"

type CameraConstructorsFactory interface {
	Register(
		ecs.ComponentType,
		func(ecs.EntityID) (Camera, error),
	)
	Build() CameraConstructors
}

type cameraConstructorsFactory struct {
	*cameraConstructors
}

func (f *cameraConstructorsFactory) Register(
	componentType ecs.ComponentType,
	ctor func(ecs.EntityID) (Camera, error),
) {
	f.constructors[componentType] = ctor
}

func (f *cameraConstructorsFactory) Build() CameraConstructors {
	return f
}

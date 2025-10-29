package cameras

import (
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) CameraConstructorsFactory {
		return &cameraConstructorsFactory{
			constructors: make(map[ecs.ComponentType]func(ecs.World) func(ecs.EntityID) (Camera, error)),
		}
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[CameraConstructors] {
		return ioc.Get[CameraConstructorsFactory](c)
	})
}

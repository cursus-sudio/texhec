package scenespkg

import (
	scenesys "frontend/modules/scenes"
	"frontend/modules/scenes/internal"
	"frontend/services/scenes"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) scenesys.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ecs.RegisterSystems(w,
				internal.NewChangeSceneSystem(ioc.Get[scenes.SceneManager](c)),
			)
			return nil
		})
	})
}

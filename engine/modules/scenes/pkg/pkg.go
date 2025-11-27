package scenespkg

import (
	scenesys "engine/modules/scenes"
	"engine/modules/scenes/internal"
	"engine/services/ecs"
	"engine/services/scenes"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) scenesys.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ecs.RegisterSystems(w,
				internal.NewChangeSceneSystem(ioc.Get[scenes.SceneManager](c)),
			)
			return nil
		})
	})
}

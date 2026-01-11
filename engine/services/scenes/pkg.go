package scenes

import (
	"engine/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneManagerBuilder {
		return NewSceneManagerBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneManager {
		return ioc.Get[SceneManagerBuilder](c).Build()
	})

	ioc.WrapService(b, func(c ioc.Dic, b events.Builder) {
		events.ListenToAll(b, func(a any) {
			m := ioc.Get[SceneManager](c)
			events.EmitAny(m.CurrentSceneWorld().Events(), a)
		})
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) ecs.World {
		return ioc.Get[SceneManager](c).CurrentSceneWorld()
	})
}

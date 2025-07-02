package scenes

import (
	"frontend/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) ecs.World {
		return ioc.Get[SceneManager](c).CurrentScene().World()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneManager {
		return newSceneManager()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneBuilder {
		return newSceneBuilder()
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b events.Builder) events.Builder {
		sceneManager := ioc.Get[SceneManager](c)
		events.ListenToAll(b, func(e any) {
			events.EmitAny(sceneManager.CurrentScene().Events(), e)
		})
		return b
	})
}

package scenes

import (
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneManagerBuilder {
		return NewSceneManagerBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneManager {
		return ioc.Get[SceneManagerBuilder](c).Build()
	})
	ioc.RegisterDependency[SceneManager, SceneManagerBuilder](b)

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b events.Builder) events.Builder {
		events.ListenToAll(b, func(a any) {
			m := ioc.Get[SceneManager](c)
			events.EmitAny(m.CurrentSceneCtx().Events, a)
		})
		return b
	})
	ioc.RegisterDependency[events.Builder, SceneManager](b)
}

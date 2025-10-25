package gamescenes

import (
	menuscene "core/scenes/menu"
	"frontend/services/scenes"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		b.AddScene(ioc.Get[menuscene.Builder](c).
			Build(menuscene.ID))
		b.MakeActive(menuscene.ID)
		return b
	})
}

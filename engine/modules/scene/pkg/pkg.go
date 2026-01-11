package scenepkg

import (
	"engine/modules/scene"
	"engine/modules/scene/internal"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// events
			Register(scene.ChangeSceneEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) scene.Service {
		return internal.NewService(c)
	})
}

package audiopkg

import (
	"frontend/modules/audio"
	"frontend/modules/audio/internal"
	"frontend/services/assets"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.Service {
		return internal.NewService(
			ioc.Get[assets.Assets](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.PlayerService { return ioc.Get[internal.Service](c) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.VolumeService { return ioc.Get[internal.Service](c) })

	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.System {
		return internal.NewSystem(ioc.Get[internal.Service](c))
	})
}

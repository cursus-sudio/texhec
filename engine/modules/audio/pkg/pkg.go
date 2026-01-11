package audiopkg

import (
	"engine/modules/audio"
	"engine/modules/audio/internal"
	"engine/services/assets"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// events
			Register(audio.StopEvent{}).
			Register(audio.PlayEvent{}).
			Register(audio.QueueEvent{}).
			Register(audio.QueueEndlessEvent{}).
			Register(audio.SetMasterVolumeEvent{}).
			Register(audio.SetChannelVolumeEvent{})
	})
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

package audiopkg

import (
	"engine/modules/audio"
	"engine/modules/audio/internal"
	"engine/services/assets"
	"engine/services/codec"
	"os"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/mix"
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
		return internal.NewService(c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.PlayerService { return ioc.Get[internal.Service](c) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.VolumeService { return ioc.Get[internal.Service](c) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.Service { return ioc.Get[internal.Service](c) })

	ioc.RegisterSingleton(b, func(c ioc.Dic) audio.System {
		return internal.NewSystem(c)
	})

	ioc.WrapService(b, func(c ioc.Dic, b assets.AssetsStorageBuilder) {
		b.RegisterExtension("wav", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			chunk, err := mix.QuickLoadWAV(source)
			if err != nil {
				return nil, err
			}
			audio := audio.NewAudioAsset(chunk, source)
			return audio, nil
		})
	})
}

package settingspkg

import (
	gameassets "core/assets"
	"core/modules/settings"
	"core/modules/settings/internal"
	"core/modules/ui"
	"engine"
	"engine/services/assets"
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
			Register(settings.EnterSettingsEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) settings.System {
		system := internal.NewSystem(
			ioc.Get[assets.Assets](c),
			ioc.Get[gameassets.GameAssets](c),
			ioc.GetServices[engine.World](c),
			ioc.Get[ui.Service](c),
		)
		return system
	})
}

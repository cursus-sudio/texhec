package settingspkg

import (
	gameassets "core/assets"
	"core/modules/settings"
	"core/modules/settings/internal"
	"core/modules/ui"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// events
			Register(settings.EnterSettingsEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) settings.System {
		system := internal.NewSystem(
			ioc.Get[assets.Assets](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[gameassets.GameAssets](c),
			ioc.Get[ecs.ToolFactory[transform.Tool]](c),
			ioc.Get[ecs.ToolFactory[render.Tool]](c),
			ioc.Get[ecs.ToolFactory[ui.Tool]](c),
			ioc.Get[ecs.ToolFactory[text.Tool]](c),
		)
		return system
	})
}

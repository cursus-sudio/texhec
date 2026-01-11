package settingspkg

import (
	"core/modules/settings"
	"core/modules/settings/internal"
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
		return internal.NewSystem(c)
	})
}

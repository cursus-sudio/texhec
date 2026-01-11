package dragpkg

import (
	"engine/modules/camera"
	"engine/modules/drag"
	"engine/modules/drag/internal"
	"engine/modules/transform"
	"engine/services/codec"

	"github.com/ogiusek/events"
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
			Register(drag.DraggableEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) drag.System {
		return internal.NewSystem(
			ioc.Get[events.Builder](c),
			ioc.Get[transform.Service](c),
			ioc.Get[camera.Service](c),
		)
	})
}

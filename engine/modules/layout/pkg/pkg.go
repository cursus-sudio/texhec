package layoutpkg

import (
	"engine/modules/layout"
	"engine/modules/layout/internal/service"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) layout.Service {
		return service.NewLayoutService(c)
	})
}

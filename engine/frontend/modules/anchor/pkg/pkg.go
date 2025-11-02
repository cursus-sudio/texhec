package anchorpkg

import (
	"frontend/modules/anchor"
	"frontend/modules/anchor/internal/anchorsys"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) anchor.System {
		return anchorsys.NewAnchorSystem(ioc.Get[logger.Logger](c))
	})
}

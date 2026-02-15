package assetspkg

import (
	"engine/modules/assets"
	"engine/modules/assets/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	parentDirectory string
}

func Package(parentDirectory string) ioc.Pkg {
	return pkg{parentDirectory}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) assets.Extensions {
		return internal.NewExtensions(c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) assets.Service {
		return internal.NewService(c, pkg.parentDirectory)
	})
}

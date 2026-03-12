package assetspkg

import (
	"engine/modules/assets"
	"engine/modules/assets/internal"
	"engine/modules/registry"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	parentDirectory string
}

func Package(parentDirectory string) ioc.Pkg {
	if len(parentDirectory) != 0 && parentDirectory[len(parentDirectory)-1] != '/' {
		parentDirectory += "/"
	}
	return pkg{parentDirectory}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) assets.Service {
		return internal.NewService(c)
	})
	ioc.WrapService(b, func(c ioc.Dic, registry registry.Service) {
		registry.Register("path", func(entity ecs.EntityID, structTagValue string) {
			assetsService := ioc.Get[assets.Service](c)
			path := assets.NewPath(pkg.parentDirectory + structTagValue)
			assetsService.Path().Set(entity, path)
		})
	})
}

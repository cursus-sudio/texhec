package worldtexture

import (
	"frontend/services/assets"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) RegisterFactory {
		return &registerFactory{
			ioc.Get[assets.AssetsStorage](c),
		}
	})
}

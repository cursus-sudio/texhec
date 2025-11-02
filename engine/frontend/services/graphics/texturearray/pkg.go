package texturearray

import (
	"frontend/services/assets"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Factory {
		return &factory{
			ioc.Get[assets.AssetsStorage](c),
		}
	})
}

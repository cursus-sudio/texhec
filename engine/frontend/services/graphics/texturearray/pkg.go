package texturearray

import (
	"frontend/services/assets"
	"shared/services/datastructures"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Factory {
		return &factory{
			ioc.Get[assets.AssetsStorage](c),
			datastructures.NewSparseArray[uint32, assets.AssetID](),
		}
	})
}

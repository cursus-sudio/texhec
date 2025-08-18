package assets

import (
	"fmt"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) AssetsStorageBuilder {
		return NewAssetsStorageBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) AssetsStorage {
		s, errs := ioc.Get[AssetsStorageBuilder](c).Build()
		if len(errs) != 0 {
			panic(fmt.Sprintf("%v\n", errs))
		}
		return s
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) AssetsCache {
		return NewCachedAssets()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) Assets {
		return &assets{
			assetStorage: ioc.Get[AssetsStorage](c),
			cachedAssets: ioc.Get[AssetsCache](c),
		}
	})
}

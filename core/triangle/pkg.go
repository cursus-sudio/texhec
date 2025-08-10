package triangle

import (
	_ "embed"
	"frontend/engine/materials/texturematerial"
	"frontend/services/assets"
	"github.com/ogiusek/ioc/v2"
	appruntime "shared/services/runtime"
)

//go:embed square.png
var textureSource []byte

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	registerAssets(b)

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			assets.Release(
				MeshAssetID,
				TextureAssetID,
				texturematerial.TextureMaterial,
			)
		})
		return b
	})
}

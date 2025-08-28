package triangle

import (
	_ "embed"
	"frontend/services/assets"
	"frontend/services/scenes"
	appruntime "shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
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
			scene := ioc.Get[scenes.SceneManager](c).CurrentSceneCtx()
			scene.World.Release()

			assets.Release(
				MeshAssetID,
				TextureAssetID,
			)
		})
		return b
	})
}

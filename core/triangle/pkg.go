package triangle

import (
	"frontend/services/assets"
	"frontend/services/scenes"
	appruntime "shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
)

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
				Texture1AssetID,
				Texture2AssetID,
				Texture3AssetID,
				Texture4AssetID,
			)
		})
		return b
	})
}

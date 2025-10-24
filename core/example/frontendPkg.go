package example

import (
	"frontend/services/assets"
	"frontend/services/scenes"
	"shared/services/ecs"
	appruntime "shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct {
}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	registerAssets(b)
	AddSceneOne(b)

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		getPersistedComponentsArray := []func(ecs.ComponentsStorage) ecs.AnyComponentArray{
			func(c ecs.ComponentsStorage) ecs.AnyComponentArray { return ecs.GetComponentsArray[someComponent](c) },
		}
		b.OnSceneLoad(func(loadedScene scenes.SceneCtx) {
			sceneManager := ioc.Get[scenes.SceneManager](c)
			unloadedScene := sceneManager.CurrentSceneCtx()
			if unloadedScene == nil {
				return
			}
			for _, getPersistedComponentsArray := range getPersistedComponentsArray {
				newArr := getPersistedComponentsArray(loadedScene.World.Components())
				oldArr := getPersistedComponentsArray(unloadedScene.World.Components())
				for _, entity := range oldArr.GetEntities() {
					loadedScene.World.EnsureEntityExists(entity)
				}
				transaction := newArr.AnyTransaction()
				for _, entity := range oldArr.GetEntities() {
					component, _ := oldArr.GetAnyComponent(entity)
					transaction.SaveAnyComponent(entity, component)
				}
				transaction.Flush()
			}
		})
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		scene1Builder := ioc.Get[SceneOneBuilder](c)
		scene1 := scene1Builder.Build(scene1Id)
		b.AddScene(scene1)

		b.MakeActive(scene1Id)
		return b
	})

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

package example

import (
	"frontend/services/scenes"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct {
}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	AddSceneOne(b)
	AddSceneTwo(b)
	AddSceneThree(b)

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

		scene2Builder := ioc.Get[SceneTwoBuilder](c)
		scene2 := scene2Builder.Build(scene2Id)
		b.AddScene(scene2)

		scene3Builder := ioc.Get[SceneThreeBuilder](c)
		scene3 := scene3Builder.Build(scene3Id)
		b.AddScene(scene3)

		b.MakeActive(scene1Id)
		return b
	})
}

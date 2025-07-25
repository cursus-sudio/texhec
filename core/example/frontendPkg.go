package example

import (
	"core/triangle"
	"frontend/engine/systems/render"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/scenes"
	"shared/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct {
}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

var scene1Id = scenes.NewSceneId("main scene")

type SceneOneBuilder scenes.SceneBuilder
type SceneOneWorld ecs.World

func AddSceneOne(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneBuilder { return scenes.NewSceneBuilder() })
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneWorld { return ecs.NewWorld() })
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, sceneBuilder SceneOneBuilder) SceneOneBuilder {
		sceneBuilder.OnLoad(func(sceneManager scenes.SceneManager, s scenes.Scene, b events.Builder) {
			world := ecs.World(ioc.Get[SceneOneWorld](c))
			triangle.AddToWorld(c, world, b)
			console := ioc.Get[console.Console](c)

			for i := 0; i < 1; i++ {
				entity := world.NewEntity()
				world.SaveComponent(entity, newSomeComponent())
			}

			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				console,
			)
			toggleSystem := NewToggledSystem(sceneManager, world, scene2Id)

			renderSystem := render.NewRenderSystem(
				world,
				ioc.Get[assets.Assets](c),
				ioc.Get[logger.Logger](c),
			)

			events.Listen(b, func(e frames.FrameEvent) {
				someSystem.Update(e)
				toggleSystem.Update(e)
				renderSystem.Update(e)
			})
		})
		return sceneBuilder
	})
}

//

var scene2Id = scenes.NewSceneId("main scene 2")

type SceneTwoBuilder scenes.SceneBuilder
type SceneTwoWorld ecs.World

func AddSceneTwo(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneTwoBuilder { return scenes.NewSceneBuilder() })
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneTwoWorld { return ecs.NewWorld() })
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, sceneBuilder SceneTwoBuilder) SceneTwoBuilder {
		sceneBuilder.OnLoad(func(sceneManager scenes.SceneManager, s scenes.Scene, b events.Builder) {
			world := ecs.World(ioc.Get[SceneTwoWorld](c))
			triangle.AddToWorld(c, world, b)

			for i := 0; i < 2; i++ {
				entity := world.NewEntity()
				world.SaveComponent(entity, newSomeComponent())
			}

			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			renderSystem := render.NewRenderSystem(
				world,
				ioc.Get[assets.Assets](c),
				ioc.Get[logger.Logger](c),
			)
			events.Listen(b, func(e frames.FrameEvent) {
				someSystem.Update(e)
				renderSystem.Update(e)
			})
		})

		return sceneBuilder
	})
}

func (FrontendPkg) Register(b ioc.Builder) {
	AddSceneOne(b)
	AddSceneTwo(b)
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		scene1Builder := ioc.Get[SceneOneBuilder](c)
		scene1 := scene1Builder.Build(scene1Id)
		b.AddScene(scene1)

		scene2Builder := ioc.Get[SceneTwoBuilder](c)
		scene2 := scene2Builder.Build(scene2Id)
		b.AddScene(scene2)

		b.MakeActive(scene1Id)
		return b
	})
}

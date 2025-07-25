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

type SceneOneBuilder events.Builder
type SceneOneWorld ecs.World

type SceneTwoBuilder events.Builder
type SceneTwoWorld ecs.World

func (FrontendPkg) Register(b ioc.Builder) {
	scene1Id := scenes.NewSceneId("main scene")
	scene2Id := scenes.NewSceneId("main scene 2")

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneBuilder { return SceneOneBuilder(events.NewBuilder()) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneWorld { return ecs.NewWorld() })

	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneTwoBuilder { return SceneTwoBuilder(events.NewBuilder()) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneTwoWorld { return ecs.NewWorld() })

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		scene := newMainScene(scene1Id, func(sceneManager scenes.SceneManager) events.Events {
			eventsBuilder := events.Builder(ioc.Get[SceneOneBuilder](c))
			world := ecs.World(ioc.Get[SceneOneWorld](c))
			triangle.AddToWorld(c, world, eventsBuilder)
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

			events.Listen(eventsBuilder, func(e frames.FrameEvent) {
				someSystem.Update(e)
				toggleSystem.Update(e)
				renderSystem.Update(e)

			})
			return eventsBuilder.Build()
		})
		b.AddScene(scene)
		b.MakeActive(scene1Id)
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		scene := newMainScene(scene2Id, func(sceneManager scenes.SceneManager) events.Events {
			eventsBuilder := events.Builder(ioc.Get[SceneTwoBuilder](c))
			world := ecs.World(ioc.Get[SceneTwoWorld](c))
			triangle.AddToWorld(c, world, eventsBuilder)

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
			events.Listen(eventsBuilder, func(e frames.FrameEvent) {
				someSystem.Update(e)
				renderSystem.Update(e)
			})
			return eventsBuilder.Build()
		})
		b.AddScene(scene)
		return b
	})
}

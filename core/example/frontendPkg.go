package example

import (
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct {
}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	scene1Id := scenes.NewSceneId("main scene")
	scene2Id := scenes.NewSceneId("main scene 2")
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		c = c.Scope(frontendscopes.Scene)

		scene := newMainScene(scene1Id, func(sceneManager scenes.SceneManager) events.Events {
			eventsBuilder := events.NewBuilder()
			world := ioc.Get[ecs.World](c)
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

			events.Listen(eventsBuilder, func(e frames.FrameEvent) {
				someSystem.Update(e)
				toggleSystem.Update(e)
			})
			return eventsBuilder.Build()
		})
		b.AddScene(scene)
		b.MakeActive(scene1Id)
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		c = c.Scope(frontendscopes.Scene)

		scene := newMainScene(scene2Id, func(sceneManager scenes.SceneManager) events.Events {
			eventsBuilder := events.NewBuilder()
			world := ioc.Get[ecs.World](c)

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
			events.Listen(eventsBuilder, func(e frames.FrameEvent) {
				someSystem.Update(e)
			})
			return eventsBuilder.Build()
		})
		b.AddScene(scene)
		return b
	})
}

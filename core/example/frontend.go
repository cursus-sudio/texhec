package example

import (
	"core/triangle"
	"fmt"
	inputssystem "frontend/engine/systems/inputs"
	"frontend/engine/systems/render"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media/inputs"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"shared/services/logger"
	"shared/services/runtime"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func AddShared(c ioc.Dic, b scenes.SceneBuilder) {
	b.OnLoad(func(sm scenes.SceneManager, s scenes.Scene, b events.Builder) {
		quitSytem := inputssystem.NewQuitSystem(
			ioc.Get[runtime.Runtime](c),
		)

		events.Listen(b, func(e sdl.QuitEvent) {
			quitSytem.Listen(e)
		})
	})

	b.OnLoad(func(sm scenes.SceneManager, s scenes.Scene, b events.Builder) {
		events.Listen(b, func(e sdl.KeyboardEvent) {
			ioc.Get[logger.Logger](c).Info(fmt.Sprintf("keyboard event is %v\nkey is %v\n", e, e.Keysym))
		})

		events.Listen(b, func(e sdl.KeyboardEvent) {
			if e.Keysym.Sym == sdl.K_q {
				ioc.Get[logger.Logger](c).Info("quiting program due to pressing 'Q'")
				ioc.Get[runtime.Runtime](c).Stop()
			}
			ioc.Get[logger.Logger](c).Info(fmt.Sprintf("keyboard event is %v\nkey is %v\n", e, e.Keysym.Sym))
		})
	})
}

type SharedSystems struct {
	inputsSystem inputssystem.InputsSystem
	renderSystem render.RenderSystem
	flushSystem  render.FlushSystem
}

func NewSharedDomain(c ioc.Dic, world ecs.World, b events.Builder) SharedSystems {
	triangle.AddToWorld(c, world, b)
	for i := 0; i < 1; i++ {
		entity := world.NewEntity()
		world.SaveComponent(entity, newSomeComponent())
	}
	return SharedSystems{
		inputsSystem: inputssystem.NewInputsSystem(ioc.Get[inputs.Api](c)),
		renderSystem: render.NewRenderSystem(world, ioc.Get[assets.Assets](c), ioc.Get[logger.Logger](c)),
		flushSystem:  render.NewFlushSystem(ioc.Get[window.Api](c)),
	}
}

func (s SharedSystems) BeforeDomain(args frames.FrameEvent) {
	s.inputsSystem.Update(args)
}
func (s SharedSystems) AfterDomain(args frames.FrameEvent) {
	s.renderSystem.Update(args)
	s.flushSystem.Update(args)
}

var scene1Id = scenes.NewSceneId("main scene")

type SceneOneBuilder scenes.SceneBuilder
type SceneOneWorld ecs.World

func AddSceneOne(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) SceneOneBuilder { return scenes.NewSceneBuilder() })
	ioc.RegisterTransient(b, func(c ioc.Dic) SceneOneWorld { return ecs.NewWorld() })
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, sceneBuilder SceneOneBuilder) SceneOneBuilder {
		AddShared(c, sceneBuilder)
		sceneBuilder.OnLoad(func(sceneManager scenes.SceneManager, s scenes.Scene, b events.Builder) {
			world := ecs.World(ioc.Get[SceneOneWorld](c))
			sharedSystems := NewSharedDomain(c, world, b)

			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			toggleSystem := NewToggledSystem(sceneManager, world, scene2Id, time.Second)

			events.Listen(b, func(e frames.FrameEvent) {
				sharedSystems.BeforeDomain(e)
				someSystem.Update(e)
				toggleSystem.Update(e)
				sharedSystems.AfterDomain(e)
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
	ioc.RegisterTransient(b, func(c ioc.Dic) SceneTwoBuilder { return scenes.NewSceneBuilder() })
	ioc.RegisterTransient(b, func(c ioc.Dic) SceneTwoWorld { return ecs.NewWorld() })
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, sceneBuilder SceneTwoBuilder) SceneTwoBuilder {
		AddShared(c, sceneBuilder)
		sceneBuilder.OnLoad(func(sceneManager scenes.SceneManager, s scenes.Scene, b events.Builder) {
			var world ecs.World = ioc.Get[SceneTwoWorld](c)
			sharedSystems := NewSharedDomain(c, world, b)

			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			// toggleSystem := NewToggledSystem(sceneManager, world, scene1Id, time.Second*3)

			events.Listen(b, func(e frames.FrameEvent) {
				sharedSystems.BeforeDomain(e)
				someSystem.Update(e)
				// toggleSystem.Update(e)
				sharedSystems.AfterDomain(e)
			})
		})

		return sceneBuilder
	})
}

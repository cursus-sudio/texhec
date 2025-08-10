package example

import (
	"core/triangle"
	"fmt"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	inputssystem "frontend/engine/systems/inputs"
	"frontend/engine/systems/mergedsystems"
	"frontend/engine/systems/mouseray"
	"frontend/engine/systems/render"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/colliders"
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
	quitSytem := inputssystem.NewQuitSystem(
		ioc.Get[runtime.Runtime](c),
	)

	b.OnLoad(func(sm scenes.SceneManager, s scenes.Scene, b events.Builder) {
		events.Listen(b, func(e sdl.QuitEvent) {
			quitSytem.Listen(e)
		})
	})

	b.OnLoad(func(sm scenes.SceneManager, s scenes.Scene, b events.Builder) {
		logger := ioc.Get[logger.Logger](c)
		events.Listen(b, func(e sdl.KeyboardEvent) {
			logger.Info(fmt.Sprintf("keyboard event is %v; key is %v", e, e.Keysym.Sym))
			if e.Keysym.Sym == sdl.K_q {
				logger.Info("quiting program due to pressing 'Q'")
				ioc.Get[runtime.Runtime](c).Stop()
			}
			if e.Keysym.Sym == sdl.K_ESCAPE {
				logger.Info("quiting program due to pressing 'ESC'")
				ioc.Get[runtime.Runtime](c).Stop()
			}
			if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_f {
				logger.Info("toggling screen size due to pressing 'F'")
				window := ioc.Get[window.Api](c)
				flags := window.Window().GetFlags()
				if flags&sdl.WINDOW_FULLSCREEN_DESKTOP == sdl.WINDOW_FULLSCREEN_DESKTOP {
					window.Window().SetFullscreen(0)
				} else {
					window.Window().SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
				}
			}
		})
	})
}

type SharedSystems struct {
	inputsSystem inputssystem.InputsSystem
	renderSystem render.RenderSystem
	flushSystem  render.FlushSystem
}

func NewSharedDomain(
	c ioc.Dic,
	world ecs.World,
	b events.Builder,
	frameSystem mergedsystems.MergedSystems[frames.FrameEvent],
) SharedSystems {
	triangle.AddToWorld(c, world, b, frameSystem)
	for i := 0; i < 1; i++ {
		entity := world.NewEntity()
		world.SaveComponent(entity, newSomeComponent())
	}

	{
		cameraRaySystem := mouseray.NewCameraRaySystem(
			world,
			ioc.Get[colliders.ColliderService](c),
			ioc.Get[window.Api](c),
			b.Events(),
			[]ecs.ComponentType{ecs.GetComponentType(projection.Ortho{}), ecs.GetComponentType(projection.Perspective{})},
			[]ecs.ComponentType{ecs.GetComponentType(mouse.MouseEvents{})},
		)
		shootRaySystem := mergedsystems.NewMergedSystems[mouseray.ShootRayEvent](func(err error) {
			ioc.Get[logger.Logger](c).Error(err)
		})
		shootRaySystem.AddSystems(cameraRaySystem.Listen)
		events.Listen(b, shootRaySystem.Listen)

		hoverSystem := inputssystem.NewHoverSystem(world, b.Events())
		events.Listen(b, hoverSystem.Listen)

		clickSystem := inputssystem.NewClickSystem(world, b.Events())
		events.Listen(b, clickSystem.Listen)

		events.Listen(b, func(event sdl.MouseMotionEvent) {
			events.Emit(b.Events(), mouseray.NewShootRayEvent())
		})
		events.Listen(b, func(event sdl.KeyboardEvent) {
			events.Emit(b.Events(), mouseray.NewShootRayEvent())
		})
		events.Listen(b, inputssystem.NewResizeSystem().Listen)
	}
	return SharedSystems{
		inputsSystem: inputssystem.NewInputsSystem(ioc.Get[inputs.Api](c)),
		renderSystem: render.NewRenderSystem(world, ioc.Get[assets.Assets](c)),
		flushSystem:  render.NewFlushSystem(ioc.Get[window.Api](c)),
	}
}

func (s SharedSystems) BeforeDomain(args frames.FrameEvent) error {
	s.inputsSystem.Update(args)
	return nil
}
func (s SharedSystems) AfterDomain(args frames.FrameEvent) error {
	if err := s.renderSystem.Update(args); err != nil {
		return err
	}
	s.flushSystem.Update(args)
	return nil
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
			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			toggleSystem := NewToggledSystem(sceneManager, world, scene2Id, time.Second)

			frameSystem := mergedsystems.NewMergedSystems[frames.FrameEvent](func(err error) {
				ioc.Get[logger.Logger](c).Error(err)
			})
			sharedSystems := NewSharedDomain(c, world, b, frameSystem)
			frameSystem.AddSystems(
				sharedSystems.BeforeDomain,
				someSystem.Update,
				toggleSystem.Update,
				sharedSystems.AfterDomain,
			)
			events.Listen(b, frameSystem.Listen)
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
			frameSystem := mergedsystems.NewMergedSystems[frames.FrameEvent](func(err error) {
				ioc.Get[logger.Logger](c).Error(err)
			})

			someSystem := NewSomeSystem(
				sceneManager,
				world,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			sharedSystems := NewSharedDomain(c, world, b, frameSystem)
			// toggleSystem := NewToggledSystem(sceneManager, world, scene1Id, time.Second*3)
			frameSystem.AddSystems(
				sharedSystems.BeforeDomain,
				someSystem.Update,
				sharedSystems.AfterDomain,
			)

			events.Listen(b, frameSystem.Listen)
		})

		return sceneBuilder
	})
}

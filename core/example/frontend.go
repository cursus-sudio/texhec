package example

import (
	"core/triangle"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/systems/inputs"
	mousesystem "frontend/engine/systems/mouse"
	"frontend/engine/systems/projections"
	"frontend/engine/systems/render"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/colliders"
	"frontend/services/console"
	"frontend/services/ecs"
	inputsmedia "frontend/services/media/inputs"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"shared/services/logger"
	"shared/services/runtime"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func AddShared[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadFirst, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.GlobalErrHandler(ctx.EventsBuilder, func(err error) { ioc.Get[logger.Logger](c).Error(err) })
		})
		return b
	})
	ioc.WrapService(b, scenes.LoadBeforeDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			quitSytem := inputs.NewQuitSystem(
				ioc.Get[runtime.Runtime](c),
			)
			events.Listen(ctx.EventsBuilder, func(e sdl.QuitEvent) {
				quitSytem.Listen(e)
			})
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			logger := ioc.Get[logger.Logger](c)

			events.Listen(ctx.EventsBuilder, func(e sdl.KeyboardEvent) {
				// logger.Info(fmt.Sprintf("keyboard event is %v; key is %v", e, e.Keysym.Sym))
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
		return b
	})

	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			for i := 0; i < 1; i++ {
				entity := ctx.World.NewEntity()
				ctx.World.SaveComponent(entity, newSomeComponent())
			}
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadBeforeDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			cameraRaySystem := mousesystem.NewCameraRaySystem(
				ctx.World,
				ioc.Get[colliders.ColliderService](c),
				ioc.Get[window.Api](c),
				ctx.Events,
				[]ecs.ComponentType{ecs.GetComponentType(projection.Ortho{}), ecs.GetComponentType(projection.Perspective{})},
				[]ecs.ComponentType{ecs.GetComponentType(mouse.MouseEvents{})},
			)
			events.Listen(ctx.EventsBuilder, func(event mousesystem.ShootRayEvent) {
				return
				if err := cameraRaySystem.Listen(event); err != nil {
					ioc.Get[logger.Logger](c).Error(err)
				}
			})

			hoverSystem := mousesystem.NewHoverSystem(ctx.World, ctx.Events)
			events.Listen(ctx.EventsBuilder, hoverSystem.Listen)

			clickSystem := mousesystem.NewClickSystem(ctx.World, ctx.Events, sdl.RELEASED)
			events.Listen(ctx.EventsBuilder, clickSystem.Listen)

			events.Listen(ctx.EventsBuilder, func(event sdl.MouseMotionEvent) {
				events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
			})
			events.Listen(ctx.EventsBuilder, func(event sdl.KeyboardEvent) {
				events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
			})
			events.Listen(ctx.EventsBuilder, inputs.NewResizeSystem().Listen)

			resizeCameraSystem := projections.NewUpdateProjectionsSystem(ctx.World, ioc.Get[window.Api](c))
			events.Listen(ctx.EventsBuilder, resizeCameraSystem.Listen)
			events.Listen(ctx.EventsBuilder, func(e sdl.WindowEvent) {
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					events.Emit(ctx.Events, projections.NewUpdateProjectionsEvent())
				}
			})
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadBeforeDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			inputsSystem := inputs.NewInputsSystem(ioc.Get[inputsmedia.Api](c))
			events.Listen(ctx.EventsBuilder, inputsSystem.Listen)
		})
		return b
	})

	triangle.AddToWorld[SceneBuilder](b)

	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			someSystem := NewSomeSystem(
				ioc.Get[scenes.SceneManager](c),
				ctx.World,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			events.ListenE(ctx.EventsBuilder, someSystem.Listen)
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadAfterDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			renderSystem := render.NewRenderSystem(ctx.World, ioc.Get[assets.Assets](c))
			events.ListenE(ctx.EventsBuilder, renderSystem.Listen)

			flushSystem := render.NewFlushSystem(ioc.Get[window.Api](c))
			events.Listen(ctx.EventsBuilder, flushSystem.Listen)
		})
		return b
	})
}

var scene1Id = scenes.NewSceneId("main scene")

type SceneOneBuilder scenes.SceneBuilder

func AddSceneOne(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneBuilder { return scenes.NewSceneBuilder() })
	AddShared[SceneOneBuilder](b)
	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, sceneBuilder SceneOneBuilder) SceneOneBuilder {
		// sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
		// 	toggleSystem := NewToggledSystem(ioc.Get[scenes.SceneManager](c), ctx.World, scene2Id, time.Second)
		// 	events.ListenE(ctx.EventsBuilder, toggleSystem.Update)
		// })
		return sceneBuilder
	})
}

//

var scene2Id = scenes.NewSceneId("main scene 2")

type SceneTwoBuilder scenes.SceneBuilder

func AddSceneTwo(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneTwoBuilder { return scenes.NewSceneBuilder() })
	AddShared[SceneTwoBuilder](b)
	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, sceneBuilder SceneTwoBuilder) SceneTwoBuilder {
		sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
			toggleSystem := NewToggledSystem(ioc.Get[scenes.SceneManager](c), ctx.World, scene1Id, time.Second)
			events.ListenE(ctx.EventsBuilder, toggleSystem.Listen)
		})
		return sceneBuilder
	})
}

var scene3Id = scenes.NewSceneId("main scene 3")

type SceneThreeBuilder scenes.SceneBuilder

func AddSceneThree(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneThreeBuilder { return scenes.NewSceneBuilder() })
	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, sceneBuilder SceneThreeBuilder) SceneThreeBuilder {
		sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
			quitSytem := inputs.NewQuitSystem(
				ioc.Get[runtime.Runtime](c),
			)
			events.Listen(ctx.EventsBuilder, func(e sdl.QuitEvent) {
				quitSytem.Listen(e)
			})
		})
		sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
			for i := 0; i < 1000000; i++ {
				ctx.World.NewEntity()
			}
		})
		sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
			inputsSystem := inputs.NewInputsSystem(ioc.Get[inputsmedia.Api](c))
			events.Listen(ctx.EventsBuilder, inputsSystem.Listen)
		})
		sceneBuilder.OnLoad(func(ctx scenes.SceneCtx) {
			someSystem := NewSomeSystem(
				ioc.Get[scenes.SceneManager](c),
				ctx.World,
				ioc.Get[backendconnection.Backend](c).Connection(),
				ioc.Get[console.Console](c),
			)
			events.ListenE(ctx.EventsBuilder, someSystem.Listen)
		})
		return sceneBuilder
	})
}

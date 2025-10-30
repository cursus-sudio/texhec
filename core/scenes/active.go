package gamescenes

import (
	"core/src/fpslogger"
	"core/src/tile"
	"frontend/engine/systems/anchor"
	"frontend/engine/systems/collider"
	genericrenderersys "frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/inputs"
	mobilecamerasystem "frontend/engine/systems/mobilecamera"
	mousesystem "frontend/engine/systems/mouse"
	"frontend/engine/systems/projections"
	quitsys "frontend/engine/systems/quit"
	"frontend/engine/systems/render"
	"frontend/engine/systems/scenes"
	textsys "frontend/engine/systems/text"
	"frontend/engine/systems/transform"
	"frontend/engine/tools/broadcollision"
	"frontend/engine/tools/cameras"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/frames"
	"frontend/services/graphics/vao/vbo"
	inputsmedia "frontend/services/media/inputs"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	MenuID     = scenes.NewSceneId("menu")
	GameID     = scenes.NewSceneId("game")
	SettingsID = scenes.NewSceneId("settings")
	CreditsID  = scenes.NewSceneId("credits")
)

type CoreSystems func(scenes.SceneCtx)

type MenuBuilder scenes.SceneBuilder
type GameBuilder scenes.SceneBuilder
type SettingsBuilder scenes.SceneBuilder
type CreditsBuilder scenes.SceneBuilder

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func AddDefaults[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadConfig, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		logger := ioc.Get[logger.Logger](c)
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.GlobalErrHandler(ctx.EventsBuilder(), func(err error) {
				logger.Error(err)
			})
		})
		return b
	})
	ioc.WrapService(b, scenes.LoadSystems, func(c ioc.Dic, s SceneBuilder) SceneBuilder {
		s.OnLoad(ioc.Get[CoreSystems](c))
		return s
	})
	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, s SceneBuilder) SceneBuilder {
		s.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events(), projectionssys.NewUpdateProjectionsEvent())
			events.Emit(ctx.Events(), mousesystem.NewShootRayEvent())
		})
		return s
	})
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		b.AddScene(ioc.Get[MenuBuilder](c).Build(MenuID))
		b.AddScene(ioc.Get[GameBuilder](c).Build(GameID))
		b.AddScene(ioc.Get[SettingsBuilder](c).Build(SettingsID))
		b.AddScene(ioc.Get[CreditsBuilder](c).Build(CreditsID))
		b.MakeActive(MenuID)
		return b
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) CoreSystems {
		return func(ctx scenes.SceneCtx) {
			logger := ioc.Get[logger.Logger](c)

			ecs.RegisterSystems(ctx,
				anchorsys.NewAnchorSystem(logger),
				transformsys.NewPivotPointSystem(logger),

				collidersys.NewColliderSystem(ioc.Get[ecs.ToolFactory[broadcollision.CollisionService]](c)),

				// mouse systems
				mousesystem.NewCameraRaySystem(
					ioc.Get[ecs.ToolFactory[broadcollision.CollisionService]](c),
					ioc.Get[window.Api](c),
					ioc.Get[ecs.ToolFactory[cameras.CameraConstructors]](c),
				),
				mousesystem.NewHoverSystem(),
				mousesystem.NewHoverEventsSystem(),
				mousesystem.NewClickSystem(logger),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(frames.FrameEvent) {
						events.Emit(w.Events(), mousesystem.NewShootRayEvent())
					})
					return nil
				}),

				// inputs systems
				inputssys.NewResizeSystem(),
				inputssys.NewInputsSystem(ioc.Get[inputsmedia.Api](c)),
				quitsys.NewQuitSystem(ioc.Get[runtime.Runtime](c)),

				// render
				rendersys.NewClearSystem(),
				rendersys.NewRenderSystem(
					ioc.Get[window.Api](c),
					2,
				),
				genericrenderersys.NewSystem(
					ioc.Get[window.Api](c),
					ioc.Get[assets.AssetsStorage](c),
					logger,
					ioc.Get[vbo.VBOFactory[genericrenderersys.Vertex]](c),
					ioc.Get[ecs.ToolFactory[cameras.CameraConstructors]](c),
					[]ecs.ComponentType{},
				),
				ioc.Get[textsys.TextRendererRegister](c),
				ioc.Get[tile.TileRenderSystemRegister](c),

				// projection and camera systems
				projectionssys.NewUpdateProjectionsSystem(ioc.Get[window.Api](c), logger),
				mobilecamerasystem.NewScrollSystem(
					logger,
					ioc.Get[ecs.ToolFactory[cameras.CameraConstructors]](c),
					ioc.Get[window.Api](c),
					0.1, 5, // min and max zoom
				),
				mobilecamerasystem.NewDragSystem(
					sdl.BUTTON_LEFT,
					ioc.Get[ecs.ToolFactory[cameras.CameraConstructors]](c),
					ioc.Get[window.Api](c),
					logger,
				),
				mobilecamerasystem.NewWasdSystem(
					ioc.Get[ecs.ToolFactory[cameras.CameraConstructors]](c),
					1.0, // speed
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(sdl.QuitEvent) {
						events.Emit(ctx.Events(), quitsys.NewQuitEvent())
					})
					events.Listen(w.EventsBuilder(), func(e sdl.WindowEvent) {
						if e.Event == sdl.WINDOWEVENT_RESIZED {
							events.Emit(ctx.Events(), projectionssys.NewUpdateProjectionsEvent())
						}
					})
					return nil
				}),

				//
				scenessys.NewChangeSceneSystem(ioc.Get[scenes.SceneManager](c)),

				// domain systems
				fpslogger.NewFpsLoggerSystem(
					ioc.Get[scenes.SceneManager](c),
					ioc.Get[console.Console](c),
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(e sdl.KeyboardEvent) {
						if e.Keysym.Sym == sdl.K_q {
							logger.Info("quiting program due to pressing 'Q'")
							events.Emit(ctx.Events(), quitsys.NewQuitEvent())
						}
						if e.Keysym.Sym == sdl.K_ESCAPE {
							logger.Info("quiting program due to pressing 'ESC'")
							events.Emit(ctx.Events(), quitsys.NewQuitEvent())
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
					return nil
				}),
			)
		}
	})
}

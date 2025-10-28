package gamescenes

import (
	"core/src/domain"
	"core/src/logs"
	"core/src/tile"
	"frontend/engine/systems/anchor"
	"frontend/engine/systems/collider"
	genericrenderersys "frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/inputs"
	mobilecamerasystem "frontend/engine/systems/mobilecamera"
	mousesystem "frontend/engine/systems/mouse"
	"frontend/engine/systems/projections"
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

type CoreSystems func(scenes.SceneCtx) []ecs.SystemRegister

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
			events.GlobalErrHandler(ctx.EventsBuilder, func(err error) {
				logger.Error(err)
			})
		})
		return b
	})
	ioc.WrapService(b, scenes.LoadSystems, func(c ioc.Dic, s SceneBuilder) SceneBuilder {
		s.OnLoad(func(ctx scenes.SceneCtx) {
			coreSystems := ioc.Get[CoreSystems](c)(ctx)
			ecs.RegisterSystems(ctx.EventsBuilder, coreSystems...)
		})
		return s
	})

	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, s SceneBuilder) SceneBuilder {
		s.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
			events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
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
		return func(ctx scenes.SceneCtx) []ecs.SystemRegister {
			logger := ioc.Get[logger.Logger](c)

			// systems

			textRenderer, err := ioc.Get[textsys.TextRendererFactory](c).New(ctx.World)
			if err != nil {
				logger.Error(err)
			}

			genericRenderer, err := genericrenderersys.NewSystem(
				ctx.World,
				ioc.Get[window.Api](c),
				ioc.Get[assets.AssetsStorage](c),
				logger,
				ioc.Get[vbo.VBOFactory[genericrenderersys.Vertex]](c),
				ioc.Get[cameras.CameraConstructorsFactory](c).Build(ctx.World),
				[]ecs.ComponentType{},
			)
			if err != nil {
				logger.Error(err)
			}

			tileService, err := ioc.Get[tile.TileRenderSystemFactory](c).
				NewSystem(ctx.World)
			if err != nil {
				logger.Error(err)
			}

			ecs.RegisterSystems(ctx.EventsBuilder,
				// core systems
				anchorsys.NewAnchorSystem(ctx.World, logger),
				transformsys.NewPivotPointSystem(ctx.World, logger),
				collidersys.NewColliderSystem(ctx.World, ioc.Get[broadcollision.CollisionServiceFactory](c)),

				// mouse systems
				mousesystem.NewCameraRaySystem(
					ctx.World,
					ioc.Get[broadcollision.CollisionServiceFactory](c)(ctx.World),
					ioc.Get[window.Api](c),
					ctx.Events,
					ioc.Get[cameras.CameraConstructorsFactory](c).Build(ctx.World),
				),
				mousesystem.NewHoverSystem(ctx.World, ctx.Events),
				mousesystem.NewHoverEventsSystem(ctx.World, ctx.Events),
				mousesystem.NewClickSystem(ctx.World, ctx.Events, sdl.RELEASED),
				ecs.NewSystemRegister(func(b events.Builder) {
					events.Listen(b, func(event frames.FrameEvent) {
						events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
					})
				}),

				// inputs systems
				inputssys.NewResizeSystem(),
				inputssys.NewQuitSystem(ioc.Get[runtime.Runtime](c)),
				inputssys.NewInputsSystem(ioc.Get[inputsmedia.Api](c)),

				// render
				rendersys.NewClearSystem(),
				rendersys.NewRenderSystem(
					ctx.World,
					ctx.Events,
					ioc.Get[window.Api](c),
					2,
				),
				genericRenderer,
				textRenderer,
				tileService,

				// projection and camera systems
				projectionssys.NewUpdateProjectionsSystem(ctx.World, ioc.Get[window.Api](c), logger),
				mobilecamerasystem.NewScrollSystem(
					ctx.World,
					logger,
					ioc.Get[cameras.CameraConstructorsFactory](c).Build(ctx.World),
					ioc.Get[window.Api](c),
					0.1, 5,
				),
				mobilecamerasystem.NewDragSystem(
					sdl.BUTTON_LEFT,
					ctx.World,
					ioc.Get[cameras.CameraConstructorsFactory](c).Build(ctx.World),
					ioc.Get[window.Api](c),
					logger,
				),
				mobilecamerasystem.NewWasdSystem(
					ctx.World,
					ioc.Get[cameras.CameraConstructorsFactory](c).Build(ctx.World),
					1.0,
				),
				ecs.NewSystemRegister(func(b events.Builder) {
					events.Listen(b, func(e sdl.WindowEvent) {
						if e.Event == sdl.WINDOWEVENT_RESIZED {
							events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
						}
					})
				}),
				scenessys.NewChangeSceneSystem(ioc.Get[scenes.SceneManager](c)),

				// domain systems
				logs.NewLogsSystem(
					ioc.Get[scenes.SceneManager](c),
					ctx.World,
					ioc.Get[console.Console](c),
				),

				ecs.NewSystemRegister(func(b events.Builder) { // some key buttons
					// export some of these to universal controlls (mainly 'F')
					events.Listen(b, func(e sdl.KeyboardEvent) {
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
				}),
				domain.NewSys(ctx.World,
					logger,
					ioc.Get[runtime.Runtime](c),
					ioc.Get[console.Console](c),
				),
			)
			return []ecs.SystemRegister{}
		}
	})
}

package gamescenes

import (
	"core/modules/fpslogger"
	"core/modules/settings"
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/modules/audio"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/drag"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/relation"
	"engine/modules/render"
	scenesys "engine/modules/scenes"
	"engine/modules/text"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"engine/services/scenes"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	MenuID       = scenes.NewSceneId("menu")
	GameID       = scenes.NewSceneId("game")
	GameClientID = scenes.NewSceneId("game client")
	SettingsID   = scenes.NewSceneId("settings")
	CreditsID    = scenes.NewSceneId("credits")
)

const (
	EffectChannel audio.Channel = iota
	MusicChannel
)

type CoreSystems func(ecs.World)

type MenuBuilder scenes.SceneBuilder
type GameBuilder scenes.SceneBuilder
type GameClientBuilder scenes.SceneBuilder
type SettingsBuilder scenes.SceneBuilder
type CreditsBuilder scenes.SceneBuilder

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func AddDefaults[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadConfig, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		logger := ioc.Get[logger.Logger](c)
		b.OnLoad(func(world ecs.World) {
			events.GlobalErrHandler(world.EventsBuilder(), func(err error) {
				logger.Warn(err)
			})
		})
		return b
	})
	ioc.WrapService(b, scenes.LoadSystems, func(c ioc.Dic, s SceneBuilder) SceneBuilder {
		s.OnLoad(ioc.Get[CoreSystems](c))
		return s
	})
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b scenes.SceneManagerBuilder) scenes.SceneManagerBuilder {
		b.AddScene(ioc.Get[MenuBuilder](c).Build(MenuID))
		b.AddScene(ioc.Get[GameBuilder](c).Build(GameID))
		b.AddScene(ioc.Get[GameClientBuilder](c).Build(GameClientID))
		b.AddScene(ioc.Get[SettingsBuilder](c).Build(SettingsID))
		b.AddScene(ioc.Get[CreditsBuilder](c).Build(CreditsID))
		b.MakeActive(MenuID)
		return b
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) CoreSystems {
		return func(world ecs.World) {
			logger := ioc.Get[logger.Logger](c)
			posFactory := ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[tile.PosComponent]]](c)
			colliderFactory := ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[tile.ColliderPos]]](c)

			temporaryInlineSystems := ecs.NewSystemRegister(func(w ecs.World) error {
				posFactory.Build(w)
				colliderFactory.Build(w)
				events.Listen(w.EventsBuilder(), func(e sdl.KeyboardEvent) {
					if e.Keysym.Sym == sdl.K_q {
						logger.Info("quiting program due to pressing 'Q'")
						events.Emit(world.Events(), inputs.NewQuitEvent())
					}
					if e.Keysym.Sym == sdl.K_ESCAPE {
						logger.Info("quiting program due to pressing 'ESC'")
						events.Emit(world.Events(), inputs.NewQuitEvent())
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
			})

			ecs.RegisterSystems(world,
				ioc.Get[netsync.StartSystem](c),
				// update {

				// inputs
				ioc.Get[inputs.System](c),

				// update
				ioc.Get[animation.System](c),
				ioc.Get[camera.System](c),
				ioc.Get[collider.System](c),
				ioc.Get[drag.System](c),
				ioc.Get[groups.System](c),
				ioc.Get[connection.System](c),
				temporaryInlineSystems,

				ioc.Get[tile.System](c),

				// ui update
				ioc.Get[ui.System](c),
				ioc.Get[settings.System](c),
				// } (update)
				ioc.Get[netsync.StopSystem](c),

				// audio
				ioc.Get[audio.System](c),

				// render
				ioc.Get[render.System](c),
				ioc.Get[tile.SystemRenderer](c),
				ioc.Get[genericrenderer.System](c),
				ioc.Get[text.System](c),
				ioc.Get[fpslogger.System](c),

				// after everything change scene
				ioc.Get[scenesys.System](c),
			)
		}
	})
}

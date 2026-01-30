package gamescenes

import (
	"core/modules/fpslogger"
	"core/modules/generation"
	"core/modules/settings"
	"core/modules/tile"
	"core/modules/ui"
	"engine"
	"engine/modules/audio"
	"engine/modules/batcher"
	"engine/modules/camera"
	"engine/modules/connection"
	"engine/modules/drag"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/render"
	"engine/modules/scene"
	"engine/modules/smooth"
	"engine/modules/text"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	MenuID       = scene.NewSceneId("menu")
	GameID       = scene.NewSceneId("game")
	GameClientID = scene.NewSceneId("game client")
	SettingsID   = scene.NewSceneId("settings")
	CreditsID    = scene.NewSceneId("credits")
)

const (
	EffectChannel audio.Channel = iota
	MusicChannel
)

type World struct {
	engine.World `inject:"1"`

	// game
	Tile       tile.Service       `inject:"1"`
	Generation generation.Service `inject:"1"`
	Ui         ui.Service         `inject:"1"`
}

type MenuBuilder scene.Scene
type GameBuilder scene.Scene
type GameClientBuilder scene.Scene
type SettingsBuilder scene.Scene
type CreditsBuilder scene.Scene

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b scene.Service) {
		b.SetScene(MenuID, scene.Scene(ioc.Get[MenuBuilder](c)))
		b.SetScene(GameID, scene.Scene(ioc.Get[GameBuilder](c)))
		b.SetScene(GameClientID, scene.Scene(ioc.Get[GameClientBuilder](c)))
		b.SetScene(SettingsID, scene.Scene(ioc.Get[SettingsBuilder](c)))
		b.SetScene(CreditsID, scene.Scene(ioc.Get[CreditsBuilder](c)))
	})

	ioc.WrapService(b, func(c ioc.Dic, b runtime.Builder) {
		b.BeforeStart(func(r runtime.Runtime) {
			logger := ioc.Get[logger.Logger](c)
			eventsBuilder := ioc.Get[events.Builder](c)
			events.GlobalErrHandler(eventsBuilder, func(err error) {
				logger.Warn(err)
			})

			temporaryInlineSystems := ecs.NewSystemRegister(func() error {
				events.Listen(eventsBuilder, func(e sdl.KeyboardEvent) {
					if e.Keysym.Sym == sdl.K_q {
						logger.Info("quiting program due to pressing 'Q'")
						events.Emit(eventsBuilder.Events(), inputs.NewQuitEvent())
					}
					if e.Keysym.Sym == sdl.K_ESCAPE {
						logger.Info("quiting program due to pressing 'ESC'")
						events.Emit(eventsBuilder.Events(), inputs.NewQuitEvent())
					}
					if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_f {
						logger.Info("toggling screen size due to pressing 'F'")
						window := ioc.Get[window.Api](c)
						flags := window.Window().GetFlags()
						if flags&sdl.WINDOW_FULLSCREEN_DESKTOP == sdl.WINDOW_FULLSCREEN_DESKTOP {
							_ = window.Window().SetFullscreen(0)
						} else {
							_ = window.Window().SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
						}
					}
				})

				return nil
			})

			errs := ecs.RegisterSystems(
				ioc.Get[netsync.StartSystem](c),
				ioc.Get[smooth.StartSystem](c),
				// update {
				ioc.Get[connection.System](c),

				// inputs
				ioc.Get[inputs.System](c),

				// update
				ioc.Get[camera.System](c),
				ioc.Get[drag.System](c),
				ioc.Get[transition.System](c),
				temporaryInlineSystems,

				ioc.Get[tile.System](c),

				// ui update
				ioc.Get[ui.System](c),
				ioc.Get[settings.System](c),
				// } (update)
				ioc.Get[smooth.StopSystem](c),
				ioc.Get[netsync.StopSystem](c),
				ioc.Get[batcher.System](c),

				// audio
				ioc.Get[audio.System](c),
				ioc.Get[inputs.ShutdownSystem](c), // after batcher and before render system

				// render
				ioc.Get[render.System](c),
				ioc.Get[tile.SystemRenderer](c),
				ioc.Get[render.SystemRenderer](c),
				ioc.Get[text.System](c),
				ioc.Get[fpslogger.System](c),
			)
			for _, err := range errs {
				logger.Warn(err)
			}
		})
	})
}

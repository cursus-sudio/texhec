package gamescenes

import (
	"core/modules/fpslogger"
	"core/modules/tile"
	"frontend/modules/anchor"
	"frontend/modules/audio"
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/drag"
	"frontend/modules/genericrenderer"
	"frontend/modules/inputs"
	"frontend/modules/render"
	scenesys "frontend/modules/scenes"
	"frontend/modules/text"
	"frontend/modules/transform"
	"frontend/services/assets"
	"frontend/services/frames"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"shared/services/ecs"
	"shared/services/logger"
	"time"

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

const (
	EffectChannel audio.Channel = iota
	MusicChannel
)

type CoreSystems func(scenes.SceneCtx)

type MenuBuilder scenes.SceneBuilder
type GameBuilder scenes.SceneBuilder
type SettingsBuilder scenes.SceneBuilder
type CreditsBuilder scenes.SceneBuilder

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
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
}

func (pkg) Register(b ioc.Builder) {
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
				// inputs
				ioc.Get[inputs.System](c),

				// update
				ioc.Get[anchor.System](c),
				ioc.Get[transform.System](c),
				ioc.Get[collider.System](c),
				ioc.Get[drag.System](c),
				ioc.Get[camera.System](c),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(e sdl.KeyboardEvent) {
						if e.Keysym.Sym == sdl.K_q {
							logger.Info("quiting program due to pressing 'Q'")
							events.Emit(ctx.Events(), inputs.NewQuitEvent())
						}
						if e.Keysym.Sym == sdl.K_ESCAPE {
							logger.Info("quiting program due to pressing 'ESC'")
							events.Emit(ctx.Events(), inputs.NewQuitEvent())
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

					// temporary inline system for animating everything
					{
						textureComponentsArray := ecs.GetComponentsArray[render.TextureComponent](w.Components())
						textureTransaction := textureComponentsArray.Transaction()
						assetsService := ioc.Get[assets.AssetsStorage](c)
						var timeElapsed time.Duration
						frameDuration := time.Millisecond * 200
						events.Listen(w.EventsBuilder(), func(e frames.FrameEvent) {
							timeElapsed += e.Delta
							if timeElapsed < frameDuration {
								return
							}
							timeElapsed -= frameDuration
							entities := textureComponentsArray.GetEntities()
							for _, entity := range entities {
								comp, err := textureComponentsArray.GetComponent(entity)
								if err != nil {
									continue
								}
								asset, err := assets.StorageGet[render.TextureAsset](assetsService, comp.Asset)
								if err != nil {
									continue
								}
								comp.Frame += 1
								comp.Frame = comp.Frame % len(asset.Images())
								textureTransaction.SaveComponent(entity, comp)
							}
							textureTransaction.Flush()
						})
					}

					return nil
				}),

				// audio
				ioc.Get[audio.System](c),

				// render
				ioc.Get[render.System](c),
				ioc.Get[genericrenderer.System](c),
				ioc.Get[tile.System](c),
				ioc.Get[text.System](c),
				ioc.Get[fpslogger.System](c),

				// after everything change scene
				ioc.Get[scenesys.System](c),
			)
		}
	})
}

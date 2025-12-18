package gamescenes

import (
	"core/modules/definition"
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
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/render"
	scenesys "engine/modules/scenes"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/uuid"
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

type World interface {
	// engine
	ecs.World
	animation.AnimationTool
	camera.CameraTool
	collider.ColliderTool
	connection.ConnectionTool
	genericrenderer.GenericRendererTool
	groups.GroupsTool
	hierarchy.HierarchyTool
	netsync.NetSyncTool
	inputs.InputsTool
	render.RenderTool
	text.TextTool
	transform.TransformTool
	uuid.UUIDTool

	// game
	definition.DefinitionTool
	tile.TileTool
	ui.UiTool
}

type world struct {
	// engine
	ecs.World
	animation.AnimationTool
	camera.CameraTool
	collider.ColliderTool
	connection.ConnectionTool
	genericrenderer.GenericRendererTool
	groups.GroupsTool
	hierarchy.HierarchyTool
	netsync.NetSyncTool
	inputs.InputsTool
	render.RenderTool
	text.TextTool
	transform.TransformTool
	uuid.UUIDTool

	// game
	definition.DefinitionTool
	tile.TileTool
	ui.UiTool
}

type WorldResolver func(ecs.World) World

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

	ioc.RegisterSingleton(b, func(c ioc.Dic) WorldResolver {
		return func(w ecs.World) World {
			world := world{World: w}
			world.AnimationTool = ioc.Get[ecs.ToolFactory[animation.World, animation.AnimationTool]](c).Build(world)
			world.UUIDTool = ioc.Get[ecs.ToolFactory[uuid.World, uuid.UUIDTool]](c).Build(world)

			world.HierarchyTool = ioc.Get[ecs.ToolFactory[hierarchy.World, hierarchy.HierarchyTool]](c).Build(world)
			world.GroupsTool = ioc.Get[ecs.ToolFactory[groups.World, groups.GroupsTool]](c).Build(world)
			world.TransformTool = ioc.Get[ecs.ToolFactory[transform.World, transform.TransformTool]](c).Build(world)

			world.ConnectionTool = ioc.Get[ecs.ToolFactory[connection.World, connection.ConnectionTool]](c).Build(world)
			world.NetSyncTool = ioc.Get[ecs.ToolFactory[netsync.World, netsync.NetSyncTool]](c).Build(world)

			world.CameraTool = ioc.Get[ecs.ToolFactory[camera.World, camera.CameraTool]](c).Build(world)
			world.ColliderTool = ioc.Get[ecs.ToolFactory[collider.World, collider.ColliderTool]](c).Build(world)
			world.GenericRendererTool = ioc.Get[ecs.ToolFactory[genericrenderer.World, genericrenderer.GenericRendererTool]](c).Build(world)
			world.InputsTool = ioc.Get[ecs.ToolFactory[inputs.World, inputs.InputsTool]](c).Build(world)
			world.RenderTool = ioc.Get[ecs.ToolFactory[render.World, render.RenderTool]](c).Build(world)
			world.TextTool = ioc.Get[ecs.ToolFactory[text.World, text.TextTool]](c).Build(world)

			world.DefinitionTool = ioc.Get[ecs.ToolFactory[definition.World, definition.DefinitionTool]](c).Build(world)
			world.TileTool = ioc.Get[ecs.ToolFactory[tile.World, tile.TileTool]](c).Build(world)
			world.UiTool = ioc.Get[ecs.ToolFactory[ui.World, ui.UiTool]](c).Build(world)
			return world
		}
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) CoreSystems {
		return func(rawWorld ecs.World) {
			world := ioc.Get[WorldResolver](c)(rawWorld)
			logger := ioc.Get[logger.Logger](c)

			temporaryInlineSystems := ecs.NewSystemRegister(func(w ecs.World) error {
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

			errs := ecs.RegisterSystems(world,
				ioc.Get[netsync.StartSystem](c),
				// update {
				ioc.Get[connection.System](c),

				// inputs
				ioc.Get[inputs.System](c),

				// update
				ioc.Get[animation.System](c),
				ioc.Get[camera.System](c),
				ioc.Get[collider.System](c),
				ioc.Get[drag.System](c),
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
			for _, err := range errs {
				logger.Warn(err)
			}
		}
	})
}

package gamescenes

import (
	"core/modules/definition"
	"core/modules/fpslogger"
	"core/modules/settings"
	"core/modules/tile"
	"core/modules/ui"
	"engine"
	"engine/modules/audio"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/drag"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/layout"
	"engine/modules/netsync"
	"engine/modules/record"
	"engine/modules/render"
	"engine/modules/scene"
	"engine/modules/smooth"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

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

type World interface {
	engine.World

	// game
	definition.DefinitionTool
	tile.TileTool
	ui.UiTool
}

type world struct {
	// engine
	ecs.World
	camera.CameraTool
	collider.ColliderTool
	connection.ConnectionTool
	genericrenderer.GenericRendererTool
	groups.GroupsTool
	hierarchy.HierarchyTool
	netsync.NetSyncTool
	record.RecordTool
	inputs.InputsTool
	layout.LayoutTool
	render.RenderTool
	text.TextTool
	transform.TransformTool
	transition.TransitionTool
	uuid.UUIDTool

	// game
	definition.DefinitionTool
	tile.TileTool
	ui.UiTool
}

type WorldResolver func(ecs.World) World

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

	ioc.RegisterSingleton(b, func(c ioc.Dic) WorldResolver {
		return func(w ecs.World) World {
			world := &world{World: w}
			world.UUIDTool = ioc.Get[uuid.ToolFactory](c).Build(world)
			world.TransitionTool = ioc.Get[transition.ToolFactory](c).Build(world)

			world.HierarchyTool = ioc.Get[hierarchy.ToolFactory](c).Build(world)
			world.GroupsTool = ioc.Get[groups.ToolFactory](c).Build(world)
			world.TransformTool = ioc.Get[transform.ToolFactory](c).Build(world)

			world.ConnectionTool = ioc.Get[connection.ToolFactory](c).Build(world)
			world.RecordTool = ioc.Get[record.ToolFactory](c).Build(world)
			world.NetSyncTool = ioc.Get[netsync.ToolFactory](c).Build(world)

			world.CameraTool = ioc.Get[camera.ToolFactory](c).Build(world)
			world.ColliderTool = ioc.Get[collider.ToolFactory](c).Build(world)
			world.GenericRendererTool = ioc.Get[genericrenderer.ToolFactory](c).Build(world)
			world.InputsTool = ioc.Get[inputs.ToolFactory](c).Build(world)
			world.LayoutTool = ioc.Get[layout.ToolFactory](c).Build(world)
			world.RenderTool = ioc.Get[render.ToolFactory](c).Build(world)
			world.TextTool = ioc.Get[text.ToolFactory](c).Build(world)

			world.DefinitionTool = ioc.Get[definition.ToolFactory](c).Build(world)
			world.TileTool = ioc.Get[tile.ToolFactory](c).Build(world)
			world.UiTool = ioc.Get[ui.ToolFactory](c).Build(world)
			return world
		}
	})

	ioc.WrapService(b, func(c ioc.Dic, rawWorld ecs.World) {
		logger := ioc.Get[logger.Logger](c)
		events.GlobalErrHandler(rawWorld.EventsBuilder(), func(err error) {
			logger.Warn(err)
		})
		world := ioc.Get[WorldResolver](c)(rawWorld)

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
						_ = window.Window().SetFullscreen(0)
					} else {
						_ = window.Window().SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					}
				}
			})

			return nil
		})

		errs := ecs.RegisterSystems(world,
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

			// audio
			ioc.Get[audio.System](c),

			// render
			ioc.Get[render.System](c),
			ioc.Get[tile.SystemRenderer](c),
			ioc.Get[genericrenderer.System](c),
			ioc.Get[text.System](c),
			ioc.Get[fpslogger.System](c),
		)
		for _, err := range errs {
			logger.Warn(err)
		}
	})
}

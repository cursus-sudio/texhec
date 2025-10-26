package menuscene

import (
	gameassets "core/assets"
	"core/src/logs"
	"core/src/tile"
	"fmt"
	"frontend/engine/components/anchor"
	"frontend/engine/components/camera"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/text"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/anchor"
	"frontend/engine/systems/collider"
	"frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/inputs"
	mobilecamerasystem "frontend/engine/systems/mobilecamera"
	mousesystem "frontend/engine/systems/mouse"
	"frontend/engine/systems/projections"
	"frontend/engine/systems/render"
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
	"slices"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

var ID = scenes.NewSceneId("main menu")

type Builder scenes.SceneBuilder

//

type QuitEvent struct{}
type OnHoveredDomainEvent struct {
	entity   ecs.EntityID
	row, col int
}
type OnClickDomainEvent struct {
	entity   ecs.EntityID
	row, col int
}

//

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) LoadConfig(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadConfig, func(c ioc.Dic, b Builder) Builder {
		logger := ioc.Get[logger.Logger](c)

		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.GlobalErrHandler(ctx.EventsBuilder, func(err error) {
				logger.Error(err)
			})
		})
		return b
	})
}

func (Pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b Builder) Builder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			world := ctx.World
			// load objects
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world.Components(), cameraEntity, transform.NewTransform())
			ecs.SaveComponent(world.Components(), cameraEntity, projection.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), cameraEntity,
				camera.NewCamera(ecs.GetComponentType(projection.Ortho{})))

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.ForestTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderersys.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world.Components(), buttonArea, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{500, 200, 1}).Val())
			ecs.SaveComponent(world.Components(), buttonArea, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).Val())

			type Button struct {
				Text    string
				OnClick any
			}
			buttons := []Button{
				{Text: "play", OnClick: QuitEvent{}},
				{Text: "settings", OnClick: QuitEvent{}},
				{Text: "credits", OnClick: QuitEvent{}},
				{Text: "exit", OnClick: QuitEvent{}},
			}
			slices.Reverse(buttons)

			for i, button := range buttons {
				entity := world.NewEntity()
				normalizedIndex := float32(i) / (float32(len(buttons)) - 1)
				ecs.SaveComponent(world.Components(), entity, transform.NewTransform().Ptr().
					SetSize(mgl32.Vec3{500, 50, 1}).Val())
				ecs.SaveComponent(world.Components(), entity, anchor.NewParentAnchor(buttonArea).Ptr().
					SetPivotPoint(mgl32.Vec3{.5, normalizedIndex, .5}).
					Val())

				ecs.SaveComponent(world.Components(), entity, mesh.NewMesh(gameassets.SquareMesh))
				ecs.SaveComponent(world.Components(), entity, texture.NewTexture(gameassets.WaterTileTextureID))
				ecs.SaveComponent(world.Components(), entity, genericrenderersys.PipelineComponent{})

				ecs.SaveComponent(world.Components(), entity, mouse.NewMouseEvents().AddLeftClickEvents(button.OnClick))
				ecs.SaveComponent(world.Components(), entity, collider.NewCollider(gameassets.SquareColliderID))

				ecs.SaveComponent(world.Components(), entity, text.Text{Text: button.Text})
				ecs.SaveComponent(world.Components(), entity, text.TextAlign{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(world.Components(), entity, text.FontSize{FontSize: 32})
			}
		})

		return b
	})
}

func (pkg Pkg) Loadsystems(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadSystems, func(c ioc.Dic, b Builder) Builder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
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
				// renderers
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

				// domain tiles system
				ecs.NewSystemRegister(func(b events.Builder) {
					tileArray := ecs.GetComponentsArray[tile.TileComponent](ctx.World.Components())
					colliderArray := ecs.GetComponentsArray[collider.Collider](ctx.World.Components())
					mouseEventsArray := ecs.GetComponentsArray[mouse.MouseEvents](ctx.World.Components())
					onChangeOrAdd := func(ei []ecs.EntityID) {
						colliderTransaction := colliderArray.Transaction()
						mouseEventsTransaction := mouseEventsArray.Transaction()
						for _, entity := range ei {
							tile, err := tileArray.GetComponent(entity)
							if err != nil {
								continue
							}

							colliderTransaction.SaveComponent(entity, collider.NewCollider(gameassets.SquareColliderID))
							mouseEventsTransaction.SaveComponent(entity, mouse.NewMouseEvents().
								AddLeftClickEvents(OnClickDomainEvent{entity, int(tile.Pos.X), int(tile.Pos.Y)}).
								AddMouseHoverEvents(OnHoveredDomainEvent{entity, int(tile.Pos.X), int(tile.Pos.Y)}),
							)
						}
						err := ecs.FlushMany(colliderTransaction, mouseEventsTransaction)
						if err != nil {
							logger.Error(err)
						}
					}

					tileArray.OnAdd(onChangeOrAdd)
					tileArray.OnChange(onChangeOrAdd)
				}),

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
				// our events
				ecs.NewSystemRegister(func(b events.Builder) {
					events.Listen(b, func(e QuitEvent) {
						ioc.Get[runtime.Runtime](c).Stop()
					})
					events.Listen(b, func(e OnHoveredDomainEvent) {
						ioc.Get[console.Console](c).Print(
							fmt.Sprintf("damn it really is hovered %v (%d, %d)\n", e.entity, e.col, e.row),
						)
					})
					events.Listen(b, func(e OnClickDomainEvent) {
						ioc.Get[console.Console](c).PrintPermanent(
							fmt.Sprintf("damn it really is clicked %v (%d, %d)\n", e.entity, e.col, e.row),
						)
					})
				}),
			)
		})
		return b
	})
}

func (Pkg) LoadInitialEvents(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, b Builder) Builder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
			events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
		})
		return b
	})
}

func (pkg Pkg) Register(b ioc.Builder) {
	pkg.LoadConfig(b)
	pkg.LoadObjects(b)
	pkg.Loadsystems(b)
	pkg.LoadInitialEvents(b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder { return scenes.NewSceneBuilder() })
}

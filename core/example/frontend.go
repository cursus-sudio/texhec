package example

import (
	"core/systems/changetransform"
	"core/tile"
	_ "embed"
	"fmt"
	"frontend/engine/components/anchor"
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mobilecamera"
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
	"frontend/engine/tools/worldprojections"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/graphics/vao/vbo"
	inputsmedia "frontend/services/media/inputs"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand/v2"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"
)

const (
	UiGroup groups.Group = iota + 1
	GameGroup
)

type QuitEvent struct{}

type OnHoveredDomainEvent struct {
	entity   ecs.EntityID
	row, col int
}
type OnClickDomainEvent struct {
	entity   ecs.EntityID
	row, col int
}

func AddGrid[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			rand := rand.New(rand.NewPCG(2077, 7137))

			tilesArray := ecs.GetComponentsArray[tile.TileComponent](ctx.World.Components())
			tilesTransaction := tilesArray.Transaction()
			rows := 100
			cols := 100
			for i := 0; i < rows*cols; i++ {
				row := i % cols
				col := i / cols
				entity := ctx.World.NewEntity()
				tileType := tile.TileMountain

				num := rand.IntN(4)

				switch num {
				case 0:
					tileType = tile.TileMountain
				case 1:
					tileType = tile.TileForest
				case 2:
					tileType = tile.TileGround
				case 3:
					tileType = tile.TileWater
				}
				tile := tile.TileComponent{
					Pos:  tile.TilePos{X: int32(row), Y: int32(col)},
					Type: tileType,
				}

				tilesTransaction.SaveComponent(entity, tile)
			}
			if err := tilesTransaction.Flush(); err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
		})
		return b
	})
}

func AddCube[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			entity := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), entity, transform.NewTransform().Ptr().
				SetPos(mgl32.Vec3{0, 0, -300}).
				SetSize(mgl32.Vec3{100, 100, 100}).Val())
			ecs.SaveComponent(ctx.World.Components(), entity, mesh.NewMesh(MeshAssetID))
			ecs.SaveComponent(ctx.World.Components(), entity, texture.NewTexture(Texture2AssetID))
			ecs.SaveComponent(ctx.World.Components(), entity, genericrenderersys.PipelineComponent{})
			ecs.SaveComponent(ctx.World.Components(), entity, projection.NewUsedProjection[projection.Perspective]())
			ecs.SaveComponent(ctx.World.Components(), entity, changetransform.Component{})
			ecs.SaveComponent(ctx.World.Components(), entity, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
		})
		return b
	})
}

func AddUi[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			for i := 0; i < 1; i++ {
				entity := ctx.World.NewEntity()
				ecs.SaveComponent(ctx.World.Components(), entity, newSomeComponent())
			}

			camera := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), camera, transform.NewTransform())
			ecs.SaveComponent(ctx.World.Components(), camera, projection.NewDynamicOrtho(
				-1000,
				+1000,
				1,
			))
			ecs.SaveComponent(ctx.World.Components(), camera, projection.NewDynamicPerspective(
				mgl32.DegToRad(90),
				0.01,
				1000,
			))
			ecs.SaveComponent(ctx.World.Components(), camera, mobilecamera.Component{})
			ecs.SaveComponent(ctx.World.Components(), camera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())

			uiCamera := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), uiCamera, transform.NewTransform())
			ecs.SaveComponent(ctx.World.Components(), uiCamera, projection.NewDynamicOrtho(
				-1000,
				+1000,
				1,
			))
			ecs.SaveComponent(ctx.World.Components(), uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			exitBtn := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), exitBtn, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 100, 1}).Val())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 1, .5}).
				Val())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, transform.NewPivotPoint(mgl32.Vec3{1, 0, .5}))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, mesh.NewMesh(MeshAssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, texture.NewTexture(Texture4AssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, genericrenderersys.PipelineComponent{})
			ecs.SaveComponent(ctx.World.Components(), exitBtn, projection.NewUsedProjection[projection.Ortho]())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, collider.NewCollider(ColliderAssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, mouse.NewMouseEvents().
				AddLeftClickEvents(QuitEvent{}),
			)
			ecs.SaveComponent(ctx.World.Components(), exitBtn, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			{
				otherBtn := ctx.World.NewEntity()
				ecs.SaveComponent(ctx.World.Components(), otherBtn, transform.NewTransform().Ptr().
					SetPos(mgl32.Vec3{-100, -100, 0}).
					SetSize(mgl32.Vec3{100, 100, 2}).Val())

				ecs.SaveComponent(ctx.World.Components(), otherBtn, mesh.NewMesh(MeshAssetID))
				ecs.SaveComponent(ctx.World.Components(), otherBtn, texture.NewTexture(Texture4AssetID))
				// ecs.SaveComponent(ctx.World.Components(), otherBtn, genericrenderersys.PipelineComponent{})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, projection.NewUsedProjection[projection.Ortho]())
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.Text{
					Text: "1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 ",
				})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.Break{Break: text.BreakWord})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.TextAlign{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.FontSize{FontSize: 14})

				ecs.SaveComponent(ctx.World.Components(), otherBtn, projection.NewUsedProjection[projection.Ortho]())
				ecs.SaveComponent(ctx.World.Components(), otherBtn, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
			}

			{
				otherBtn := ctx.World.NewEntity()
				ecs.SaveComponent(ctx.World.Components(), otherBtn, transform.NewTransform().Ptr().
					SetPos(mgl32.Vec3{-70, -70, 50}).
					SetSize(mgl32.Vec3{100, 100, 2}).Val())

				ecs.SaveComponent(ctx.World.Components(), otherBtn, mesh.NewMesh(MeshAssetID))
				ecs.SaveComponent(ctx.World.Components(), otherBtn, texture.NewTexture(Texture3AssetID))
				ecs.SaveComponent(ctx.World.Components(), otherBtn, genericrenderersys.PipelineComponent{})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, projection.NewUsedProjection[projection.Ortho]())
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.Text{
					Text: "1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 1234 ",
				})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.Break{Break: text.BreakWord})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.TextAlign{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(ctx.World.Components(), otherBtn, text.FontSize{FontSize: 14})

				ecs.SaveComponent(ctx.World.Components(), otherBtn, projection.NewUsedProjection[projection.Ortho]())
				ecs.SaveComponent(ctx.World.Components(), otherBtn, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
			}
		})

		return b
	})
}

func AddShared[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	AddGrid[SceneBuilder](b)
	// AddCube[SceneBuilder](b)
	AddUi[SceneBuilder](b)

	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.GlobalErrHandler(ctx.EventsBuilder, func(err error) { ioc.Get[logger.Logger](c).Error(err) })

			textRenderer, err := ioc.Get[textsys.TextRendererFactory](c).New(ctx.World)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
				return
			}

			genericRenderer, err := genericrenderersys.NewSystem(
				ctx.World,
				ioc.Get[window.Api](c),
				ioc.Get[assets.AssetsStorage](c),
				ioc.Get[logger.Logger](c),
				ioc.Get[vbo.VBOFactory[genericrenderersys.Vertex]](c),
				ioc.Get[cameras.CameraConstructors](c),
				[]ecs.ComponentType{},
			)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}

			tileFactory := ioc.Get[tile.TileRenderSystemFactory](c)
			tileService, err := tileFactory.NewSystem(ctx.World)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}

			ecs.RegisterSystems(ctx.EventsBuilder,
				rendersys.NewClearSystem(),
				anchorsys.NewAnchorSystem(ctx.World, ioc.Get[logger.Logger](c)),
				transformsys.NewPivotPointSystem(ctx.World, ioc.Get[logger.Logger](c)),
				inputssys.NewQuitSystem(
					ioc.Get[runtime.Runtime](c),
				),
				collidersys.NewColliderSystem(ctx.World, ioc.Get[broadcollision.CollisionServiceFactory](c)),
				mousesystem.NewCameraRaySystem(
					ctx.World,
					ioc.Get[broadcollision.CollisionServiceFactory](c)(ctx.World),
					ioc.Get[window.Api](c),
					ctx.Events,
					ioc.Get[cameras.CameraConstructors](c),
				),
				mousesystem.NewHoverSystem(ctx.World, ctx.Events),
				mousesystem.NewHoverEventsSystem(ctx.World, ctx.Events),
				mousesystem.NewClickSystem(ctx.World, ctx.Events, sdl.RELEASED),
				inputssys.NewResizeSystem(),
				projectionssys.NewUpdateProjectionsSystem(ctx.World, ioc.Get[window.Api](c), ioc.Get[logger.Logger](c)),
				inputssys.NewInputsSystem(ioc.Get[inputsmedia.Api](c)),
				NewSomeSystem(
					ioc.Get[scenes.SceneManager](c),
					ctx.World,
					ioc.Get[backendconnection.Backend](c).Connection(),
					ioc.Get[console.Console](c),
				),
				changetransform.NewSystem(ctx.World, ioc.Get[logger.Logger](c)),
				genericRenderer,
				textRenderer,
				tileService,
				mobilecamerasystem.NewScrollSystem(
					ctx.World,
					ioc.Get[logger.Logger](c),
					ioc.Get[cameras.CameraConstructors](c),
					ioc.Get[window.Api](c),
					0.1, 5,
				),
				mobilecamerasystem.NewDragSystem(
					sdl.BUTTON_LEFT,
					ctx.World,
					ioc.Get[cameras.CameraConstructors](c),
					ioc.Get[window.Api](c),
					ioc.Get[logger.Logger](c),
				),
				mobilecamerasystem.NewWasdSystem(
					ctx.World,
					ioc.Get[cameras.CameraConstructors](c),
					1.0,
				),
				rendersys.NewRenderSystem(
					ctx.World,
					ctx.Events,
					ioc.Get[window.Api](c),
					2,
				),
			)

			events.Listen(ctx.EventsBuilder, func(e QuitEvent) {
				ioc.Get[runtime.Runtime](c).Stop()
			})

			logger := ioc.Get[logger.Logger](c)

			events.Listen(ctx.EventsBuilder, func(e sdl.KeyboardEvent) {
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

			events.Listen(ctx.EventsBuilder, func(event sdl.MouseMotionEvent) {
				events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
			})
			events.Listen(ctx.EventsBuilder, func(event sdl.KeyboardEvent) {
				events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
			})

			events.Listen(ctx.EventsBuilder, func(e sdl.WindowEvent) {
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
				}
			})

			events.Listen(ctx.EventsBuilder, func(e OnHoveredDomainEvent) {
				ioc.Get[console.Console](c).Print(
					fmt.Sprintf("damn it really is hovered %v (%d, %d)\n", e.entity, e.col, e.row),
				)
			})

			events.Listen(ctx.EventsBuilder, func(e OnClickDomainEvent) {
				ioc.Get[console.Console](c).PrintPermanent(
					fmt.Sprintf("damn it really is clicked %v (%d, %d)\n", e.entity, e.col, e.row),
				)
			})

			projectionsRegister := worldprojections.NewWorldProjectionsRegister(
				ecs.GetComponentType(projection.Ortho{}),
				ecs.GetComponentType(projection.Perspective{}),
			)
			ctx.World.SaveRegister(projectionsRegister)

			{ // add listener to add all related components with grid
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

						colliderTransaction.SaveComponent(entity, collider.NewCollider(ColliderAssetID))
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
			}

			events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
		})
		return b
	})

}

var scene1Id = scenes.NewSceneId("main scene")

type SceneOneBuilder scenes.SceneBuilder

func AddSceneOne(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) SceneOneBuilder { return scenes.NewSceneBuilder() })
	AddShared[SceneOneBuilder](b)
}

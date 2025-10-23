package triangle

import (
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
	"frontend/engine/systems/genericrenderer"
	mobilecamerasystem "frontend/engine/systems/mobilecamera"
	"frontend/engine/systems/projections"
	textsys "frontend/engine/systems/text"
	"frontend/engine/systems/transform"
	"frontend/engine/tools/cameras"
	"frontend/engine/tools/worldprojections"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"math/rand/v2"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	UiGroup groups.Group = iota + 1
	GameGroup
)

func AddToWorld[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
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

			type QuitEvent struct{}

			events.Listen(ctx.EventsBuilder, func(e QuitEvent) {
				ioc.Get[runtime.Runtime](c).Stop()
			})

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
				// AddMouseHoverEvents(QuitEvent{}).
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

				textRenderer, err := ioc.Get[textsys.TextRendererFactory](c).New(ctx.World)
				if err != nil {
					ioc.Get[logger.Logger](c).Error(err)
					return
				}
				ecs.RegisterSystems(ctx.EventsBuilder,
					textRenderer,
				)
			}
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadFirst, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			anchorsys.NewAnchorSystem(ctx.World, ioc.Get[logger.Logger](c))
			transformsys.NewPivotPointSystem(ctx.World, ioc.Get[logger.Logger](c))
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			projectionsRegister := worldprojections.NewWorldProjectionsRegister(
				ecs.GetComponentType(projection.Ortho{}),
				ecs.GetComponentType(projection.Perspective{}),
			)
			ctx.World.SaveRegister(projectionsRegister)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) { // cube
			// entity := ctx.World.NewEntity()
			// ecs.SaveComponent(ctx.World.Components(), entity, transform.NewTransform().Ptr().
			// 	SetPos(mgl32.Vec3{0, 0, -300}).
			// 	SetSize(mgl32.Vec3{100, 100, 100}).Val())
			// ecs.SaveComponent(ctx.World.Components(), entity, mesh.NewMesh(MeshAssetID))
			// ecs.SaveComponent(ctx.World.Components(), entity, texture.NewTexture(Texture2AssetID))
			// ecs.SaveComponent(ctx.World.Components(), entity, genericrenderersys.PipelineComponent{})
			// ecs.SaveComponent(ctx.World.Components(), entity, projection.NewUsedProjection[projection.Perspective]())
			// ecs.SaveComponent(ctx.World.Components(), entity, ChangeTransformOverTimeComponent{})
			// ecs.SaveComponent(ctx.World.Components(), entity, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
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
			ecs.RegisterSystems(ctx.EventsBuilder,
				genericRenderer,
				NewChangeTransformOverTimeSystem(ctx.World, ioc.Get[logger.Logger](c)),
			)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			type OnHoveredDomainEvent struct {
				entity   ecs.EntityID
				row, col int
			}
			type OnClickDomainEvent struct {
				entity   ecs.EntityID
				row, col int
			}

			{
				// add listener to add all related components with grid
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
						ioc.Get[logger.Logger](c).Error(err)
					}
				}

				tileArray.OnAdd(onChangeOrAdd)
				tileArray.OnChange(onChangeOrAdd)
			}

			{
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
			}

			{
				tileFactory := ioc.Get[tile.TileRenderSystemFactory](c)
				s, err := tileFactory.NewSystem(ctx.World)
				if err != nil {
					ioc.Get[logger.Logger](c).Error(err)
				}
				ecs.RegisterSystems(ctx.EventsBuilder,
					s,
				)
			}

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

		b.OnLoad(func(ctx scenes.SceneCtx) {
			ecs.RegisterSystems(ctx.EventsBuilder,
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
			)
		})
		return b
	})
}

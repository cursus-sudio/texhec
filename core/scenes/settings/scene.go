package settingsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
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
	"frontend/engine/systems/genericrenderer"
	mousesystem "frontend/engine/systems/mouse"
	"frontend/engine/systems/projections"
	"frontend/engine/systems/scenes"
	"frontend/services/console"
	"frontend/services/scenes"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"
	"slices"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

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
	ioc.WrapService(b, scenes.LoadConfig, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
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
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			world := ctx.World
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world.Components(), cameraEntity, transform.NewTransform())
			ecs.SaveComponent(world.Components(), cameraEntity, projection.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), cameraEntity,
				camera.NewCamera(ecs.GetComponentType(projection.Ortho{})))

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(mgl32.Vec3{1, 1, .5}))
			ecs.SaveComponent(world.Components(), signature, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 0, .5}).
				Val())

			ecs.SaveComponent(world.Components(), signature, text.Text{Text: "settings"})
			ecs.SaveComponent(world.Components(), signature, text.TextAlign{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), signature, text.FontSize{FontSize: 32})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.GroundTileTextureID))
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
				{Text: "play", OnClick: scenessys.NewChangeSceneEvent(gamescenes.GameID)},
				{Text: "settings", OnClick: scenessys.NewChangeSceneEvent(gamescenes.SettingsID)},
				{Text: "credits", OnClick: scenessys.NewChangeSceneEvent(gamescenes.CreditsID)},
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
	ioc.WrapService(b, scenes.LoadSystems, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			logger := ioc.Get[logger.Logger](c)

			// systems
			coreSystems := ioc.Get[gamescenes.CoreSystems](c)(ctx)
			ecs.RegisterSystems(ctx.EventsBuilder, coreSystems...)

			ecs.RegisterSystems(ctx.EventsBuilder,
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
	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projectionssys.NewUpdateProjectionsEvent())
			events.Emit(ctx.Events, mousesystem.NewShootRayEvent())
		})
		return b
	})
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.SettingsBuilder { return scenes.NewSceneBuilder() })

	pkg.LoadConfig(b)
	pkg.LoadObjects(b)
	pkg.Loadsystems(b)
	pkg.LoadInitialEvents(b)
}

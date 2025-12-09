package gamescene

import (
	gameassets "core/assets"
	"core/modules/definition"
	"core/modules/settings"
	"core/modules/tile"
	"core/modules/ui"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/scenes"
	"math/rand/v2"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

const (
	UiGroup groups.Group = iota + 1
	GameGroup
)

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.GameBuilder) gamescenes.GameBuilder {
		b.OnLoad(func(world scenes.SceneCtx) {
			uiCamera := world.NewEntity()
			ecs.SaveComponent(world, uiCamera, camera.NewOrtho(-1000, +1000))
			ecs.SaveComponent(world, uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
			ecs.SaveComponent(world, uiCamera, ui.UiCameraComponent{})

			gameCamera := world.NewEntity()
			ecs.SaveComponent(world, gameCamera, uuid.New([16]byte{48}))
			ecs.SaveComponent(world, gameCamera, camera.NewOrtho(-1000, +1000))
			ecs.SaveComponent(world, gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
			ecs.SaveComponent(world, gameCamera, camera.NewMobileCamera())
			ecs.SaveComponent(world, gameCamera, camera.NewCameraLimits(
				mgl32.Vec3{0, 0, -1000},                // min
				mgl32.Vec3{100 * 100, 100 * 100, 1000}, // max
			))

			signature := world.NewEntity()
			ecs.SaveComponent(world, signature, transform.NewPos(5, 5, 1))
			ecs.SaveComponent(world, signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world, signature, transform.NewPivotPoint(0, .5, 0))
			ecs.SaveComponent(world, signature, hierarchy.NewParent(uiCamera))
			ecs.SaveComponent(world, signature, transform.NewParent(transform.RelativePos))
			ecs.SaveComponent(world, signature, transform.NewParentPivotPoint(0, 0, .5))
			ecs.SaveComponent(world, signature, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world, signature, text.TextComponent{Text: "game"})
			ecs.SaveComponent(world, signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world, signature, text.BreakComponent{Break: text.BreakNone})

			settingsEntity := world.NewEntity()
			ecs.SaveComponent(world, settingsEntity, transform.NewPos(10, -10, 0))
			ecs.SaveComponent(world, settingsEntity, transform.NewSize(50, 50, 1))
			ecs.SaveComponent(world, settingsEntity, transform.NewPivotPoint(0, 1, .5))
			ecs.SaveComponent(world, settingsEntity, hierarchy.NewParent(uiCamera))
			ecs.SaveComponent(world, settingsEntity, transform.NewParent(transform.RelativePos))
			ecs.SaveComponent(world, settingsEntity, transform.NewParentPivotPoint(0, 1, .5))
			ecs.SaveComponent(world, settingsEntity, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world, settingsEntity, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, settingsEntity, render.NewTexture(gameassets.SettingsTextureID))
			ecs.SaveComponent(world, settingsEntity, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world, settingsEntity, inputs.NewMouseLeftClick(settings.EnterSettingsEvent{}))
			ecs.SaveComponent(world, settingsEntity, inputs.KeepSelectedComponent{})
			ecs.SaveComponent(world, settingsEntity, collider.NewCollider(gameassets.SquareColliderID))
			if gameassets.IsServer {
				rand := rand.New(rand.NewPCG(2077, 7137))

				tilesTypeArray := ecs.GetComponentsArray[definition.DefinitionLinkComponent](world)
				tilesPosArray := ecs.GetComponentsArray[tile.PosComponent](world)
				tilesTypeTransaction := tilesTypeArray.Transaction()
				tilesPosTransaction := tilesPosArray.Transaction()
				rows := 100
				cols := 100
				{
					unit := world.NewEntity()
					tilesPosTransaction.SaveComponent(unit, tile.NewPos(1, 1, tile.UnitLayer))
					tilesTypeTransaction.SaveComponent(unit, definition.NewLink(definition.TileU1))
				}
				for i := 0; i < rows*cols; i++ {
					row := i % cols
					col := i / cols
					entity := world.NewEntity()
					tileType := definition.TileMountain

					num := rand.IntN(4)

					switch num {
					case 0:
						tileType = definition.TileMountain
					case 1:
						tileType = definition.TileGround
					case 2:
						tileType = definition.TileForest
					case 3:
						tileType = definition.TileWater
					}
					tilesPosTransaction.SaveComponent(entity, tile.NewPos(int32(row), int32(col), tile.GroundLayer))
					tilesTypeTransaction.SaveComponent(entity, definition.NewLink(tileType))
				}

				{
					unit := world.NewEntity()
					tilesPosTransaction.SaveComponent(unit, tile.NewPos(0, 0, tile.UnitLayer))
					tilesTypeTransaction.SaveComponent(unit, definition.NewLink(definition.TileU1))
				}
				err := ecs.FlushMany(tilesTypeTransaction, tilesPosTransaction)
				ioc.Get[logger.Logger](c).Warn(err)
			}

			connectionToolFactory := ioc.Get[ecs.ToolFactory[connection.Tool]](c)
			connectionTool := connectionToolFactory.Build(world)
			if gameassets.IsServer {
				connectionTool.Host(":8080", func(cc connection.ConnectionComponent) {
					entity := world.NewEntity()
					ecs.SaveComponent(world, entity, netsync.ClientComponent{})
					ecs.SaveComponent(world, entity, cc)
				})
			} else {
				comp, err := connectionTool.Connect(":8080")
				if err != nil {
					panic("nie ma serwera")
				}
				entity := world.NewEntity()
				ecs.SaveComponent(world, entity, netsync.ServerComponent{})
				ecs.SaveComponent(world, entity, comp)
			}

		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.GameBuilder](b)

	pkg.LoadObjects(b)
}

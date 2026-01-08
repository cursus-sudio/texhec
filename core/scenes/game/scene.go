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
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/scenes"
	"errors"
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

func addScene(
	world gamescenes.World,
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
	isServer bool,
) {
	rows := 1000
	cols := 1000

	uiCamera := world.NewEntity()
	world.Camera().Ortho().Set(uiCamera, camera.NewOrtho(-1000, +1000))
	world.Groups().Component().Set(uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
	world.Ui().UiCamera().Set(uiCamera, ui.UiCameraComponent{})

	gameCamera := world.NewEntity()
	world.UUID().Component().Set(gameCamera, uuid.New([16]byte{48}))
	world.Camera().Ortho().Set(gameCamera, camera.NewOrtho(-1000, +1000))
	world.Groups().Component().Set(gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
	world.Camera().Mobile().Set(gameCamera, camera.NewMobileCamera())
	world.Camera().Limits().Set(gameCamera, camera.NewCameraLimits(
		mgl32.Vec3{0, 0, -1000},
		mgl32.Vec3{100 * float32(rows), 100 * float32(cols), 1000},
	))

	signature := world.NewEntity()
	world.Transform().Pos().Set(signature, transform.NewPos(5, 5, 1))
	world.Transform().Size().Set(signature, transform.NewSize(100, 50, 1))
	world.Transform().PivotPoint().Set(signature, transform.NewPivotPoint(0, .5, 0))
	world.Hierarchy().SetParent(signature, uiCamera)
	world.Transform().Parent().Set(signature, transform.NewParent(transform.RelativePos))
	world.Transform().ParentPivotPoint().Set(signature, transform.NewParentPivotPoint(0, 0, .5))
	world.Groups().Component().Set(signature, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

	world.Text().FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
	world.Text().Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

	settingsEntity := world.NewEntity()
	world.Transform().Pos().Set(settingsEntity, transform.NewPos(10, -10, 0))
	world.Transform().Size().Set(settingsEntity, transform.NewSize(50, 50, 1))
	world.Transform().PivotPoint().Set(settingsEntity, transform.NewPivotPoint(0, 1, .5))
	world.Hierarchy().SetParent(settingsEntity, uiCamera)
	world.Transform().Parent().Set(settingsEntity, transform.NewParent(transform.RelativePos))
	world.Transform().ParentPivotPoint().Set(settingsEntity, transform.NewParentPivotPoint(0, 1, .5))
	world.Groups().Component().Set(settingsEntity, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

	world.Render().Mesh().Set(settingsEntity, render.NewMesh(gameAssets.SquareMesh))
	world.Render().Texture().Set(settingsEntity, render.NewTexture(gameAssets.Hud.Settings))
	world.GenericRenderer().Pipeline().Set(settingsEntity, genericrenderer.PipelineComponent{})

	world.Inputs().LeftClick().Set(settingsEntity, inputs.NewLeftClick(settings.EnterSettingsEvent{}))
	world.Inputs().KeepSelected().Set(settingsEntity, inputs.KeepSelectedComponent{})
	world.Collider().Component().Set(settingsEntity, collider.NewCollider(gameAssets.SquareCollider))

	if isServer {
		rand := rand.New(rand.NewPCG(2077, 7137))

		tilesTypeArray := ecs.GetComponentsArray[definition.DefinitionLinkComponent](world)
		tilesPosArray := ecs.GetComponentsArray[tile.PosComponent](world)
		{
			unit := world.NewEntity()
			world.Tile().Pos().Set(unit, tile.NewPos(1, 1, tile.UnitLayer))
			world.Definition().Link().Set(unit, definition.NewLink(definition.TileU1))
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
			tilesPosArray.Set(entity, tile.NewPos(row, col, tile.GroundLayer))
			tilesTypeArray.Set(entity, definition.NewLink(tileType))
		}

		{
			unit := world.NewEntity()
			tilesPosArray.Set(unit, tile.NewPos(0, 0, tile.UnitLayer))
			tilesTypeArray.Set(unit, definition.NewLink(definition.TileU1))
		}
	}

	if isServer {
		err := world.Connection().Host(":8000", func(cc connection.ConnectionComponent) {
			entity := world.NewEntity()
			world.NetSync().Client().Set(entity, netsync.ClientComponent{})
			world.Connection().Component().Set(entity, cc)
		})
		logger.Warn(err)
	} else {
		comp, err := world.Connection().Connect(":8000")
		if err != nil {
			logger.Warn(errors.New("there is no server"))
		}
		entity := world.NewEntity()
		world.NetSync().Server().Set(entity, netsync.ServerComponent{})
		world.Connection().Component().Set(entity, comp)
	}
}

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.GameBuilder) gamescenes.GameBuilder {
		b.OnLoad(func(rawWorld ecs.World) {
			world := ioc.Get[gamescenes.WorldResolver](c)(rawWorld)
			addScene(
				world,
				ioc.Get[gameassets.GameAssets](c),
				ioc.Get[logger.Logger](c),
				true, // is server
			)
		})

		return b
	})
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.GameClientBuilder) gamescenes.GameClientBuilder {
		b.OnLoad(func(rawWorld ecs.World) {
			world := ioc.Get[gamescenes.WorldResolver](c)(rawWorld)
			addScene(
				world,
				ioc.Get[gameassets.GameAssets](c),
				ioc.Get[logger.Logger](c),
				false, // is server
			)
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.GameBuilder](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameClientBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.GameClientBuilder](b)

	pkg.LoadObjects(b)
}

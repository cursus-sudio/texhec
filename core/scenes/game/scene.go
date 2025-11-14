package gamescene

import (
	gameassets "core/assets"
	"core/modules/tilerenderer"
	gamescenes "core/scenes"
	"frontend/modules/anchor"
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/genericrenderer"
	"frontend/modules/groups"
	"frontend/modules/inputs"
	"frontend/modules/render"
	scenessys "frontend/modules/scenes"
	"frontend/modules/text"
	"frontend/modules/transform"
	"frontend/services/scenes"
	"math/rand/v2"
	"shared/services/ecs"
	"shared/services/logger"

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
			ecs.SaveComponent(world.Components(), uiCamera, transform.NewTransform())
			ecs.SaveComponent(world.Components(), uiCamera, camera.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			gameCamera := world.NewEntity()
			ecs.SaveComponent(world.Components(), gameCamera, transform.NewTransform())
			ecs.SaveComponent(world.Components(), gameCamera, camera.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())
			ecs.SaveComponent(world.Components(), gameCamera, camera.NewMobileCamera())
			ecs.SaveComponent(world.Components(), gameCamera, camera.NewCameraLimits(
				mgl32.Vec3{0, 0, -1000},                // min
				mgl32.Vec3{100 * 100, 100 * 100, 1000}, // max
			))

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(mgl32.Vec3{1, .5, .5}))
			ecs.SaveComponent(world.Components(), signature, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 0, .5}).
				SetOffset(mgl32.Vec3{5, 5}).
				Val())
			ecs.SaveComponent(world.Components(), signature, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world.Components(), signature, text.TextComponent{Text: "game"})
			ecs.SaveComponent(world.Components(), signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetOffset(mgl32.Vec3{0, 0, -100}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
			ecs.SaveComponent(world.Components(), background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderer.PipelineComponent{})

			quit := world.NewEntity()
			ecs.SaveComponent(world.Components(), quit, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{50, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), quit, transform.NewPivotPoint(mgl32.Vec3{0, 0, .5}))
			ecs.SaveComponent(world.Components(), quit, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{1, 1, .5}).
				SetOffset(mgl32.Vec3{-10, -10}).
				Val())
			ecs.SaveComponent(world.Components(), quit, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world.Components(), quit, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), quit, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), quit, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world.Components(), quit, inputs.NewMouseEvents().
				AddLeftClickEvents(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world.Components(), quit, inputs.KeepSelectedComponent{})
			ecs.SaveComponent(world.Components(), quit, collider.NewCollider(gameassets.SquareColliderID))

			ecs.SaveComponent(world.Components(), quit, text.TextComponent{Text: "X"})
			ecs.SaveComponent(world.Components(), quit, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), quit, text.FontSizeComponent{FontSize: 32})

			rand := rand.New(rand.NewPCG(2077, 7137))

			tilesTypeArray := ecs.GetComponentsArray[tilerenderer.TileTextureComponent](world.Components())
			tilesPosArray := ecs.GetComponentsArray[tilerenderer.TilePosComponent](world.Components())
			tilesTypeTransaction := tilesTypeArray.Transaction()
			tilesPosTransaction := tilesPosArray.Transaction()
			rows := 100
			cols := 100
			for i := 0; i < rows*cols; i++ {
				row := i % cols
				col := i / cols
				entity := world.NewEntity()
				tileType := tilerenderer.TileMountain

				num := rand.IntN(4)

				switch num {
				case 0:
					tileType = tilerenderer.TileMountain
				case 1:
					tileType = tilerenderer.TileGround
				case 2:
					tileType = tilerenderer.TileForest
				case 3:
					tileType = tilerenderer.TileWater
				}
				tilesPosTransaction.SaveAnyComponent(entity, tilerenderer.NewTilePos(int32(row), int32(col), 0))
				tilesTypeTransaction.SaveComponent(entity, tilerenderer.TileTextureComponent{Texture: tileType})
			}

			{
				unit := world.NewEntity()
				tilesPosTransaction.SaveAnyComponent(unit, tilerenderer.NewTilePos(0, 0, 0))
				tilesTypeTransaction.SaveComponent(unit, tilerenderer.TileTextureComponent{Texture: tilerenderer.TileU1})
			}
			err := ecs.FlushMany(tilesTypeTransaction, tilesPosTransaction)
			ioc.Get[logger.Logger](c).Warn(err)
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.GameBuilder](b)

	pkg.LoadObjects(b)
}

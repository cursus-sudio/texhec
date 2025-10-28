package gamescene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"core/src/tile"
	"frontend/engine/components/anchor"
	"frontend/engine/components/camera"
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mobilecamera"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/text"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/scenes"
	"frontend/services/scenes"
	"math/rand/v2"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

const (
	UiGroup groups.Group = iota + 1
	GameGroup
)

func (Pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.GameBuilder) gamescenes.GameBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			world := ctx.World
			uiCamera := world.NewEntity()
			ecs.SaveComponent(world.Components(), uiCamera, transform.NewTransform())
			ecs.SaveComponent(world.Components(), uiCamera, projection.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), uiCamera,
				camera.NewCamera(ecs.GetComponentType(projection.Ortho{})))
			ecs.SaveComponent(world.Components(), uiCamera, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			gameCamera := world.NewEntity()
			ecs.SaveComponent(world.Components(), gameCamera, transform.NewTransform())
			ecs.SaveComponent(world.Components(), gameCamera, projection.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), gameCamera,
				camera.NewCamera(ecs.GetComponentType(projection.Ortho{})))
			ecs.SaveComponent(world.Components(), gameCamera, mobilecamera.Component{})
			ecs.SaveComponent(world.Components(), gameCamera, groups.EmptyGroups().Ptr().Enable(GameGroup).Val())

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(mgl32.Vec3{1, .5, .5}))
			ecs.SaveComponent(world.Components(), signature, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 0, .5}).
				SetOffset(mgl32.Vec3{5, 5}).
				Val())
			ecs.SaveComponent(world.Components(), signature, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world.Components(), signature, text.Text{Text: "game"})
			ecs.SaveComponent(world.Components(), signature, text.FontSize{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.Break{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetOffset(mgl32.Vec3{0, 0, -100}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderersys.PipelineComponent{})
			ecs.SaveComponent(world.Components(), background, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			quit := world.NewEntity()
			ecs.SaveComponent(world.Components(), quit, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{50, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), quit, transform.NewPivotPoint(mgl32.Vec3{0, 0, .5}))
			ecs.SaveComponent(world.Components(), quit, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{1, 1, .5}).
				SetOffset(mgl32.Vec3{-10, -10}).
				Val())
			ecs.SaveComponent(world.Components(), quit, groups.EmptyGroups().Ptr().Enable(UiGroup).Val())

			ecs.SaveComponent(world.Components(), quit, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), quit, texture.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), quit, genericrenderersys.PipelineComponent{})

			ecs.SaveComponent(world.Components(), quit, mouse.NewMouseEvents().
				AddLeftClickEvents(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world.Components(), quit, collider.NewCollider(gameassets.SquareColliderID))

			ecs.SaveComponent(world.Components(), quit, text.Text{Text: "X"})
			ecs.SaveComponent(world.Components(), quit, text.TextAlign{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), quit, text.FontSize{FontSize: 32})

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

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.GameBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.GameBuilder](b)

	pkg.LoadObjects(b)
}

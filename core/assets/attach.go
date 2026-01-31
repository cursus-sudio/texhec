package gameassets

import (
	"core/modules/tile"
	"engine/modules/collider"
	"engine/modules/render"
	"engine/modules/transition"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/graphics/vao/ebo"
	"engine/services/logger"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type GameAssets struct {
	Tiles        TileAssets
	Hud          HudAssets
	ExampleAudio assets.AssetID `path:"audio.wav"`

	Blank          assets.AssetID `path:"blank texture"`
	SquareMesh     assets.AssetID `path:"square mesh"`
	SquareCollider assets.AssetID `path:"square collider"`
	FontAsset      assets.AssetID `path:"font1.ttf"`
}

type HudAssets struct {
	Btn         assets.AssetID `path:"hud/btn.png"`
	Settings    assets.AssetID `path:"hud/settings.png"`
	Background1 assets.AssetID `path:"hud/bg1.gif"`
	Background2 assets.AssetID `path:"hud/bg2.gif"`
}

type TileAssets struct {
	Grass    assets.AssetID `path:"tiles/grass.biom"`
	Sand     assets.AssetID `path:"tiles/sand.biom"`
	Mountain assets.AssetID `path:"tiles/mountain.biom"`
	Water    assets.AssetID `path:"tiles/water.biom"`

	Unit assets.AssetID `path:"tiles/u1.png"`
}

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Assets(b ioc.Builder) {
	// register specific files
	ioc.WrapService(b, func(c ioc.Dic, b assets.AssetsStorageBuilder) {
		gameAssets := ioc.Get[GameAssets](c)
		b.RegisterAsset(gameAssets.Blank, func() (any, error) {
			img := image.NewRGBA(image.Rect(0, 0, 1, 1))
			white := color.RGBA{255, 255, 255, 255}
			img.Set(0, 0, white)
			asset, err := render.NewTextureStorageAsset(img)
			return asset, err
		})
		b.RegisterAsset(gameAssets.SquareMesh, func() (any, error) {
			vertices := []render.Vertex{
				{Pos: [3]float32{1, 1, 1}, TexturePos: [2]float32{1, 1}},
				{Pos: [3]float32{1, -1, 1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{-1, -1, 1}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{-1, 1, 1}, TexturePos: [2]float32{0, 1}},
			}

			indices := []ebo.Index{
				0, 1, 2,
				0, 2, 3,
			}
			asset := render.NewMeshStorageAsset(vertices, indices)
			return asset, nil
		})

		b.RegisterAsset(gameAssets.SquareCollider, func() (any, error) {
			asset := collider.NewColliderStorageAsset(
				[]collider.AABB{collider.NewAABB(mgl32.Vec3{-1, -1}, mgl32.Vec3{1, 1})},
				[]collider.Range{collider.NewRange(collider.Leaf, 0, 2)},
				[]collider.Polygon{
					collider.NewPolygon(mgl32.Vec3{-1, -1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
					collider.NewPolygon(mgl32.Vec3{+1, +1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
				})
			return asset, nil
		})
	})

	// register assets
	ioc.RegisterSingleton(b, func(c ioc.Dic) GameAssets {
		logger := ioc.Get[logger.Logger](c)
		assetsService := ioc.Get[assets.AssetModule](c)

		gameAssets := GameAssets{}
		logger.Warn(assetsService.InitializeProperties(&gameAssets))
		return gameAssets
	})

	ioc.WrapService(b, func(c ioc.Dic, s tile.TileAssets) {
		gameAssets := ioc.Get[GameAssets](c)
		assets := datastructures.NewSparseArray[tile.Type, assets.AssetID]()
		assets.Set(tile.TileSand, gameAssets.Tiles.Sand)
		assets.Set(tile.TileMountain, gameAssets.Tiles.Mountain)
		assets.Set(tile.TileGrass, gameAssets.Tiles.Grass)
		assets.Set(tile.TileWater, gameAssets.Tiles.Water)
		s.AddType(assets)
	})
}

//
//
//

const (
	_ transition.EasingID = iota
	LinearEasingFunction
	MyEasingFunction
	EaseOutElastic
)

func (pkg) Animations(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b transition.EasingService) {
		b.Set(LinearEasingFunction, func(t transition.Progress) transition.Progress {
			return t
		})
		b.Set(MyEasingFunction, func(t transition.Progress) transition.Progress {
			const n1 = 7.5625
			const d1 = 2.75

			if t < 1/d1 { // First segment of the bounce (rising curve)
				return n1 * t * t
			} else if t < 2/d1 { // Second segment (peak of the first bounce)
				t -= 1.5 / d1
				return n1*t*t + 0.75
			} else if t < 2.5/d1 { // Third segment (peak of the second, smaller bounce)
				t -= 2.25 / d1
				return n1*t*t + 0.9375
			} else { // Final segment (settling)
				t -= 2.625 / d1
				return n1*t*t + 0.984375
			}
		})
		b.Set(EaseOutElastic, func(t transition.Progress) transition.Progress {
			const c1 float64 = 10
			const c2 float64 = .75
			const c3 float64 = (2 * math.Pi) / 3
			if t == 0 {
				return 0
			}
			if t == 1 {
				return 1
			}
			x := float64(t)
			x = math.Pow(2, -c1*x)*
				math.Sin((x*c1-c2)*c3) +
				1
			return transition.Progress(x)
		})
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	pkg.Assets(b)
	pkg.Animations(b)
}

package gameassets

import (
	"bytes"
	"core/modules/definition"
	"core/modules/tile"
	"engine/modules/animation"
	"engine/modules/audio"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	gtexture "engine/services/graphics/texture"
	"engine/services/graphics/vao/ebo"
	"engine/services/logger"
	appruntime "engine/services/runtime"
	"engine/services/scenes"
	"image"
	_ "image/png"
	"math"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/mix"
	"golang.org/x/image/font/opentype"
)

type GameAssets struct {
	Tiles        TileAssets
	Hud          HudAssets
	ExampleAudio assets.AssetID `path:"audio.wav"`

	SquareMesh     assets.AssetID `path:"square mesh"`
	SquareCollider assets.AssetID `path:"square collider"`
	FontAsset      assets.AssetID `path:"font3.ttf"`
}

type HudAssets struct {
	Btn assets.AssetID `path:"hud/btn.png"`
	// BtnAspectRatio mgl32.Vec3
	Settings assets.AssetID `path:"hud/settings.png"`
}

type TileAssets struct {
	Mountain assets.AssetID `path:"tiles/mountain.png"`
	Ground   assets.AssetID `path:"tiles/ground.png"`
	Forest   assets.AssetID `path:"tiles/forest.png"`
	Water    assets.AssetID `path:"tiles/water.png"`

	Unit assets.AssetID `path:"tiles/u1.png"`
}

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Assets(b ioc.Builder) {
	// register specific files
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterExtension("wav", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			chunk, err := mix.QuickLoadWAV(source)
			if err != nil {
				return nil, err
			}
			audio := audio.NewAudioAsset(chunk, source)
			return audio, nil
		})

		b.RegisterExtension("png", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			imgFile := bytes.NewBuffer(source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}

			img = gtexture.FlipImage(img)
			if !strings.Contains(string(id), "tiles") {
				img = TrimTransparentBackground(img)
			}
			return render.NewTextureStorageAsset(img)
		})

		b.RegisterExtension("ttf", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			font, err := opentype.Parse(source)
			if err != nil {
				return nil, err
			}
			asset := text.NewFontFaceAsset(*font)
			return asset, nil
		})

		gameAssets := ioc.Get[GameAssets](c)
		b.RegisterAsset(gameAssets.SquareMesh, func() (any, error) {
			vertices := []genericrenderer.Vertex{
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

		return b
	})

	// register assets
	ioc.RegisterSingleton(b, func(c ioc.Dic) GameAssets {
		logger := ioc.Get[logger.Logger](c)
		assetsService := ioc.Get[assets.AssetModule](c)

		gameAssets := GameAssets{}
		logger.Warn(assetsService.InitializeProperties(&gameAssets))
		return gameAssets
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			scene := ioc.Get[scenes.SceneManager](c).CurrentSceneWorld()
			scene.ReleaseGlobals()

			assets.ReleaseAll()
		})
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.TileAssets) tile.TileAssets {
		gameAssets := ioc.Get[GameAssets](c)
		assets := datastructures.NewSparseArray[definition.DefinitionID, assets.AssetID]()
		assets.Set(definition.TileMountain, gameAssets.Tiles.Mountain)
		assets.Set(definition.TileGround, gameAssets.Tiles.Ground)
		assets.Set(definition.TileForest, gameAssets.Tiles.Forest)
		assets.Set(definition.TileWater, gameAssets.Tiles.Water)
		assets.Set(definition.TileU1, gameAssets.Tiles.Unit)
		s.AddType(assets)
		return s
	})
}

//
//
//

const (
	ChangeColorsAnimation animation.AnimationID = iota

	// game scene events
	ShowMenuAnimation
	HideMenuAnimation
)
const (
	MyEasingFunction animation.EasingFunctionID = iota
	LinearEasingFunction
	EaseOutElastic
)

func (pkg) Animations(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b animation.AnimationSystemBuilder) animation.AnimationSystemBuilder {
		b.AddEasingFunction(LinearEasingFunction, func(t animation.AnimationState) animation.AnimationState { return t })
		b.AddEasingFunction(MyEasingFunction, func(t animation.AnimationState) animation.AnimationState {
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
		b.AddEasingFunction(EaseOutElastic, func(t animation.AnimationState) animation.AnimationState {
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
			return animation.AnimationState(x)
		})
		b.AddAnimation(ChangeColorsAnimation, animation.NewAnimation(
			[]animation.Event{},
			[]animation.Transition{
				animation.NewTransition(
					render.NewColor(mgl32.Vec4{1, 0, 1, 1}),
					render.NewColor(mgl32.Vec4{1, 1, 1, 1}),
					MyEasingFunction,
				),
			},
		))
		b.AddAnimation(ShowMenuAnimation, animation.NewAnimation(
			[]animation.Event{},
			[]animation.Transition{
				animation.NewTransition(
					transform.NewPivotPoint(0, 1, .5),
					transform.NewPivotPoint(1, 1, .5),
					LinearEasingFunction,
				),
			},
		))
		b.AddAnimation(HideMenuAnimation, animation.NewAnimation(
			[]animation.Event{},
			[]animation.Transition{
				animation.NewTransition(
					transform.NewPivotPoint(1, 1, .5),
					transform.NewPivotPoint(0, 1, .5),
					LinearEasingFunction,
				),
			},
		))
		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	pkg.Assets(b)
	pkg.Animations(b)
}

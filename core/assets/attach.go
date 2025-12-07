package gameassets

import (
	"bytes"
	"core/modules/definition"
	"core/modules/tile"
	_ "embed"
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
	appruntime "engine/services/runtime"
	"engine/services/scenes"
	"image"
	_ "image/png"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/mix"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// var IsServer bool = true
var IsServer bool = false

//go:embed files/1.png
var mountainSource []byte

//go:embed files/2.png
var groundSource []byte

//go:embed files/3.png
var forestSource []byte

//go:embed files/4.png
var waterSource []byte

//go:embed files/u1.png
var u1Source []byte

//go:embed files/settings.png
var settingsSource []byte

//go:embed files/audio.wav
var audioSource []byte

var fontSource []byte = goregular.TTF

const (
	SquareMesh assets.AssetID = "square mesh"

	MountainTileTextureID assets.AssetID = "mountain tile texture"
	GroundTileTextureID   assets.AssetID = "ground tile texture"
	ForestTileTextureID   assets.AssetID = "forest tile texture"
	WaterTileTextureID    assets.AssetID = "water tile texture"

	U1TextureID assets.AssetID = "u1 texture"

	SettingsTextureID assets.AssetID = "settings texture"

	SquareColliderID assets.AssetID = "square collider"
	FontAssetID      assets.AssetID = "font_asset"

	AudioID assets.AssetID = "audio.wav"
)

const (
	ChangeColorsAnimation animation.AnimationID = iota
	ButtonAnimation

	// game scene events
	ShowMenuAnimation
	HideMenuAnimation
)
const (
	MyEasingFunction animation.EasingFunctionID = iota
	LinearEasingFunction
	EaseOutElastic
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
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
		b.AddAnimation(ButtonAnimation, animation.NewAnimation(
			[]animation.Event{},
			[]animation.Transition{
				animation.NewTransition(
					render.NewTextureFrameComponent(0),
					render.NewTextureFrameComponent(.6),
					LinearEasingFunction,
				),
				animation.NewTransition(
					render.NewColor(mgl32.Vec4{1, 1, 0, 1}),
					render.NewColor(mgl32.Vec4{1, 1, 1, 1}),
					MyEasingFunction,
				).SetStart(0).SetEnd(.5),
				animation.NewTransition(
					render.NewColor(mgl32.Vec4{1, 1, 1, 1}),
					render.NewColor(mgl32.Vec4{1, 1, 0, 1}),
					MyEasingFunction,
				).SetStart(.5).SetEnd(1),
			},
		))
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
	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			scene := ioc.Get[scenes.SceneManager](c).CurrentSceneCtx()
			scene.Release()

			assets.ReleaseAll()
		})
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.TileAssets) tile.TileAssets {
		assets := datastructures.NewSparseArray[definition.DefinitionID, assets.AssetID]()
		assets.Set(definition.TileMountain, MountainTileTextureID)
		assets.Set(definition.TileGround, GroundTileTextureID)
		assets.Set(definition.TileForest, ForestTileTextureID)
		assets.Set(definition.TileWater, WaterTileTextureID)
		assets.Set(definition.TileU1, U1TextureID)
		s.AddType(assets)
		return s
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(SquareMesh, func() (any, error) {
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

		b.RegisterAsset(MountainTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(mountainSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}

			img = gtexture.FlipImage(img)
			asset := render.NewTextureStorageAsset(img)
			return asset, nil
		})

		b.RegisterAsset(GroundTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(groundSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			img = gtexture.FlipImage(img)
			asset := render.NewTextureStorageAsset(img)
			return asset, nil
		})

		b.RegisterAsset(ForestTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(forestSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			img = gtexture.FlipImage(img)
			asset := render.NewTextureStorageAsset(img)
			return asset, nil
		})

		b.RegisterAsset(WaterTileTextureID, func() (any, error) {
			img1File := bytes.NewBuffer(waterSource)
			img1, _, err := image.Decode(img1File)
			if err != nil {
				return nil, err
			}
			img1 = gtexture.FlipImage(img1)

			img2File := bytes.NewBuffer(waterSource)
			img2, _, err := image.Decode(img2File)
			if err != nil {
				return nil, err
			}
			img2 = Rotate90Clockwise(img2)
			// img2 = gtexture.FlipImage(img2)
			asset := render.NewTextureStorageAsset(img1, img2)
			// asset := render.NewTextureStorageAsset(img2, img1)
			return asset, nil
		})

		b.RegisterAsset(U1TextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(u1Source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			img = gtexture.FlipImage(img)
			asset := render.NewTextureStorageAsset(img)
			return asset, nil
		})

		b.RegisterAsset(SettingsTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(settingsSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			img = gtexture.FlipImage(img)
			asset := render.NewTextureStorageAsset(img)
			return asset, nil
		})

		b.RegisterAsset(SquareColliderID, func() (any, error) {
			asset := collider.NewColliderStorageAsset(
				[]collider.AABB{collider.NewAABB(mgl32.Vec3{-1, -1}, mgl32.Vec3{1, 1})},
				[]collider.Range{collider.NewRange(collider.Leaf, 0, 2)},
				[]collider.Polygon{
					collider.NewPolygon(mgl32.Vec3{-1, -1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
					collider.NewPolygon(mgl32.Vec3{+1, +1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
				})
			return asset, nil
		})

		b.RegisterAsset(FontAssetID, func() (any, error) {
			font, err := opentype.Parse(fontSource)
			if err != nil {
				return nil, err
			}
			asset := text.NewFontFaceAsset(*font)
			return asset, nil
		})

		b.RegisterAsset(AudioID, func() (any, error) {
			chunk, err := mix.QuickLoadWAV(audioSource)
			if err != nil {
				return nil, err
			}
			audio := audio.NewAudioAsset(chunk)
			return audio, nil
		})
		return b
	})
}

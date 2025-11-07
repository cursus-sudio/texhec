package gameassets

import (
	"bytes"
	"core/modules/tile"
	_ "embed"
	"frontend/modules/audio"
	"frontend/modules/collider"
	"frontend/modules/genericrenderer"
	"frontend/modules/render"
	"frontend/modules/text"
	"frontend/services/assets"
	gtexture "frontend/services/graphics/texture"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/scenes"
	"image"
	_ "image/png"
	"shared/services/datastructures"
	appruntime "shared/services/runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/mix"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

//go:embed files/1.png
var mountainSource []byte

//go:embed files/2.png
var groundSource []byte

//go:embed files/3.png
var forestSource []byte

//go:embed files/4.png
var waterSource []byte

//go:embed files/audio.wav
var audioSource []byte

var fontSource []byte = goregular.TTF

const (
	SquareMesh assets.AssetID = "square mesh"

	MountainTileTextureID assets.AssetID = "mountain tile texture"
	GroundTileTextureID   assets.AssetID = "ground tile texture"
	ForestTileTextureID   assets.AssetID = "forest tile texture"
	WaterTileTextureID    assets.AssetID = "water tile texture"

	SquareColliderID assets.AssetID = "square collider"
	FontAssetID      assets.AssetID = "font_asset"

	AudioID assets.AssetID = "audio.wav"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}
func Rotate90Clockwise(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	W := bounds.Dx()
	H := bounds.Dy()

	rotatedRect := image.Rect(0, 0, H, W)
	rotatedImg := image.NewRGBA(rotatedRect)

	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			c := img.At(x, y)

			newX := H - 1 - y
			newY := x

			rotatedImg.Set(newX, newY, c)
		}
	}

	return rotatedImg
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			scene := ioc.Get[scenes.SceneManager](c).CurrentSceneCtx()
			scene.Release()

			assets.ReleaseAll()
		})
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.TileTool) tile.TileTool {
		assets := datastructures.NewSparseArray[uint32, assets.AssetID]()
		assets.Set(tile.TileMountain, MountainTileTextureID)
		assets.Set(tile.TileGround, GroundTileTextureID)
		assets.Set(tile.TileForest, ForestTileTextureID)
		assets.Set(tile.TileWater, WaterTileTextureID)
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

package gameassets

import (
	"bytes"
	"core/src/tile"
	_ "embed"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/text"
	"frontend/engine/components/texture"
	"frontend/engine/systems/genericrenderer"
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

var fontSource []byte = goregular.TTF

const (
	SquareMesh assets.AssetID = "square mesh"

	MountainTileTextureID assets.AssetID = "mountain tile texture"
	GroundTileTextureID   assets.AssetID = "ground tile texture"
	ForestTileTextureID   assets.AssetID = "forest tile texture"
	WaterTileTextureID    assets.AssetID = "water tile texture"

	SquareColliderID assets.AssetID = "square collider"
	FontAssetID      assets.AssetID = "font_asset"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			scene := ioc.Get[scenes.SceneManager](c).CurrentSceneCtx()
			scene.Release()

			assets.ReleaseAll()
		})
		return b
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.TileRenderSystemFactory) tile.TileRenderSystemFactory {
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
			vertices := []genericrenderersys.Vertex{
				// Front face
				{Pos: [3]float32{-1, 1, 1}, TexturePos: [2]float32{0, 1}},
				{Pos: [3]float32{-1, -1, 1}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{1, -1, 1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{1, 1, 1}, TexturePos: [2]float32{1, 1}},

				// Back face
				{Pos: [3]float32{-1, 1, -1}, TexturePos: [2]float32{1, 1}},
				{Pos: [3]float32{-1, -1, -1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{1, -1, -1}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{1, 1, -1}, TexturePos: [2]float32{0, 1}},
			}

			indices := []ebo.Index{
				// Front face
				0, 1, 2,
				0, 2, 3,
				// Back face
				4, 5, 6,
				4, 6, 7,
				// Top face
				3, 7, 4,
				3, 4, 0,
				// Bottom face
				1, 5, 6,
				1, 6, 2,
				// Right face
				2, 6, 7,
				2, 7, 3,
				// Left face
				5, 1, 0,
				5, 0, 4,
			}
			asset := mesh.NewMeshStorageAsset(vertices, indices)
			return asset, nil
		})

		b.RegisterAsset(MountainTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(mountainSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}

			imgInverse := gtexture.FlipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(GroundTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(groundSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := gtexture.FlipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(ForestTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(forestSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := gtexture.FlipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(WaterTileTextureID, func() (any, error) {
			imgFile := bytes.NewBuffer(waterSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := gtexture.FlipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
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
		return b
	})
}

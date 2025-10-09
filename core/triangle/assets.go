package triangle

import (
	"bytes"
	"core/tile"
	_ "embed"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/texture"
	"frontend/engine/systems/genericrenderer"
	"frontend/services/assets"
	"frontend/services/graphics/vao/ebo"
	"image"
	"image/draw"
	_ "image/png"
	"shared/services/datastructures"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

//go:embed assets/1.png
var texture1Source []byte

//go:embed assets/2.png
var texture2Source []byte

//go:embed assets/3.png
var texture3Source []byte

//go:embed assets/4.png
var texture4Source []byte

const (
	MeshAssetID     assets.AssetID = "vao_asset"
	Texture1AssetID assets.AssetID = "texture1_asset"
	Texture2AssetID assets.AssetID = "texture2_asset"
	Texture3AssetID assets.AssetID = "texture3_asset"
	Texture4AssetID assets.AssetID = "texture4_asset"
	ColliderAssetID assets.AssetID = "collider_asset"
)

func flipImage(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		// Use draw.Draw to copy a row of pixels.
		// The destination rectangle is a single row at the flipped Y coordinate.
		destY := bounds.Dy() - 1 - y
		destRect := image.Rect(0, destY, bounds.Dx(), destY+1)

		// The source rectangle is a single row at the original Y coordinate.
		srcRect := image.Rect(0, y, bounds.Dx(), y+1)

		// Copy the source row to the destination row.
		draw.Draw(newImg, destRect, img, srcRect.Min, draw.Src)
	}
	return newImg
}

func registerAssets(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.TileRenderSystemFactory) tile.TileRenderSystemFactory {
		assets := datastructures.NewSparseArray[uint32, assets.AssetID]()
		assets.Set(tile.TileMountain, Texture1AssetID)
		assets.Set(tile.TileGround, Texture2AssetID)
		assets.Set(tile.TileForest, Texture3AssetID)
		assets.Set(tile.TileWater, Texture4AssetID)
		s.AddType(assets)
		return s
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(MeshAssetID, func() (any, error) {
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

		b.RegisterAsset(Texture1AssetID, func() (any, error) {
			imgFile := bytes.NewBuffer(texture1Source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := flipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(Texture2AssetID, func() (any, error) {
			imgFile := bytes.NewBuffer(texture2Source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := flipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(Texture3AssetID, func() (any, error) {
			imgFile := bytes.NewBuffer(texture3Source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := flipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(Texture4AssetID, func() (any, error) {
			imgFile := bytes.NewBuffer(texture4Source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			imgInverse := flipImage(img)
			asset := texture.NewTextureStorageAsset(imgInverse)
			return asset, nil
		})

		b.RegisterAsset(ColliderAssetID, func() (any, error) {
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

}

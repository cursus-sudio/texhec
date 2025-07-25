package triangle

import (
	"bytes"
	_ "embed"
	"frontend/components/mesh"
	"frontend/components/texture"
	"frontend/services/assets"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"

	"github.com/ogiusek/ioc/v2"
)

func registerAssets(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(MeshAssetID, func() (assets.StorageAsset, error) {
			vertices := []vbo.Vertex{
				// Front face
				{Pos: [3]float32{1, 1, 1}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{1, -1, 1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{-1, -1, 1}, TexturePos: [2]float32{1, 1}},
				{Pos: [3]float32{-1, 1, 1}, TexturePos: [2]float32{0, 1}},

				// Back face
				{Pos: [3]float32{-1, 1, -1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{-1, -1, -1}, TexturePos: [2]float32{1, 1}},
				{Pos: [3]float32{1, -1, -1}, TexturePos: [2]float32{0, 1}},
				{Pos: [3]float32{1, 1, -1}, TexturePos: [2]float32{0, 0}},
			}

			indices := []ebo.Index{
				0, 1, 2,
				0, 2, 3,

				4, 5, 6,
				4, 6, 7,

				0, 1, 6,
				0, 6, 7,

				3, 2, 5, //
				3, 5, 4,

				0, 3, 4,
				0, 4, 7, //

				1, 5, 2, //
				1, 6, 5,
			}
			asset := mesh.NewMeshStorageAsset(vertices, indices)
			return asset, nil
		})
		b.RegisterAsset(TextureAssetID, func() (assets.StorageAsset, error) {
			asset := texture.NewTextureStorageAsset(bytes.NewBuffer(textureSource))
			return asset, nil
		})
		return b
	})

}

package triangle

import (
	"bytes"
	_ "embed"
	"frontend/engine/assets/material"
	"frontend/engine/components/mesh"
	"frontend/engine/components/texture"
	"frontend/services/assets"
	"frontend/services/graphics/vao/ebo"
	"image"
	_ "image/png"

	"github.com/ogiusek/ioc/v2"
)

func registerAssets(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(MeshAssetID, func() (any, error) {
			vertices := []material.Vertex{
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
		b.RegisterAsset(TextureAssetID, func() (any, error) {
			imgFile := bytes.NewBuffer(textureSource)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}
			asset := texture.NewTextureStorageAsset(img)
			return asset, nil
		})
		return b
	})

}

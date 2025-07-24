package triangle

import (
	"bytes"
	_ "embed"
	"frontend/components/material"
	"frontend/components/mesh"
	"frontend/components/texture"
	"frontend/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/materials/texturematerial"
	appruntime "shared/services/runtime"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//go:embed square.png
var textureSource []byte

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

const (
	MeshAssetID    assets.AssetID = "vao_asset"
	TextureAssetID assets.AssetID = "texture_asset"
)

func (FrontendPkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(MeshAssetID, func() (assets.StorageAsset, error) {
			// vertices := []vbo.Vertex{
			// 	{Pos: [3]float32{-1, -1, 0}, TexturePos: [2]float32{0, 0}},
			// 	{Pos: [3]float32{1, -1, 0}, TexturePos: [2]float32{1, 0}},
			// 	{Pos: [3]float32{-1, 1, 0}, TexturePos: [2]float32{0, 1}},
			// 	{Pos: [3]float32{1, 1, 0}, TexturePos: [2]float32{1, 1}},
			// }
			// indices := []ebo.Index{
			// 	0, 1, 3,
			// 	0, 2, 3,
			// }

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
		return b
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(TextureAssetID, func() (assets.StorageAsset, error) {
			asset := texture.NewTextureStorageAsset(bytes.NewBuffer(textureSource))
			return asset, nil
		})
		return b
	})

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b events.Builder) events.Builder {
		var t time.Duration
		assetsService := ioc.Get[assets.Assets](c)
		world := ioc.Get[ecs.World](c)
		entity := world.NewEntity()

		originalTransform := transform.NewTransform().
			SetPos(transform.NewPos(0, 0, 300)).
			SetSize(transform.NewSize(100, 100, 100))
		world.SaveComponent(entity, originalTransform)
		world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
		world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial))
		world.SaveComponent(entity, texture.NewTexture(TextureAssetID))

		events.Listen(b, func(e frames.FrameEvent) {
			meshAsset, err := assets.GetAsset[mesh.MeshCachedAsset](assetsService, MeshAssetID)
			if err != nil {
				panic(err)
			}
			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](assetsService, texturematerial.TextureMaterial)
			if err != nil {
				panic(err)
			}

			t += e.Delta

			materialAsset.OnFrame(world)

			{
				transformComponent := originalTransform

				radians := mgl32.DegToRad(float32(t.Seconds()) * 100)
				// radians := mgl32.DegToRad(45)
				rotation := mgl32.QuatIdent().
					Mul(mgl32.QuatRotate(radians, mgl32.Vec3{1, 0, 0})).
					// Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 1, 0})).
					Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 0, 1}))
					// Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 1, 1}))
				transformComponent.Rotation = rotation

				// transformComponent.Size.Y *= 1 + float32(t.Seconds())
				// transformComponent.Pos.X = float32(t.Seconds()) * 100

				world.SaveComponent(entity, transformComponent)
			}

			if err := materialAsset.UseForEntity(world, entity); err != nil {
				panic(err)
			}
			meshAsset.VAO().Draw()
		})
		return b
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		assets := ioc.Get[assets.Assets](c)
		b.OnStop(func(r appruntime.Runtime) {
			assets.Release(
				MeshAssetID,
				TextureAssetID,
				texturematerial.TextureMaterial,
			)
		})
		return b
	})
}

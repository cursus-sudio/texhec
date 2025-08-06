package triangle

import (
	_ "embed"
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/materials/texturematerial"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/media/window"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

const (
	MeshAssetID    assets.AssetID = "vao_asset"
	TextureAssetID assets.AssetID = "texture_asset"
)

func AddToWorld(c ioc.Dic, world ecs.World, b events.Builder) {
	{
		width, height := ioc.Get[window.Api](c).Window().GetSize()
		camera := world.NewEntity()
		world.SaveComponent(camera, transform.NewTransform())
		world.SaveComponent(camera, projection.NewOrtho(float32(width), float32(height), -1000, 1000))

		fovY := mgl32.DegToRad(90)
		aspectRatio := float32(width) / float32(height)
		var near, far float32 = 0.001, 1000
		world.SaveComponent(camera, projection.NewPerspective(fovY, aspectRatio, near, far))
	}

	{
		entity := world.NewEntity()

		world.SaveComponent(entity, transform.NewTransform().
			SetPos(mgl32.Vec3{0, 0, 300}).
			SetSize(mgl32.Vec3{100, 100, 100}))
		world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
		world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial3D, texturematerial.TextureMaterial2D))
		world.SaveComponent(entity, texture.NewTexture(TextureAssetID))
		world.SaveComponent(entity, ChangeTransformOverTimeComponent{})
	}

	system := ChangeTransformOverTimeSystem{World: world}

	events.Listen(b, system.Update)
}

package triangle

import (
	_ "embed"
	"frontend/components/material"
	"frontend/components/mesh"
	"frontend/components/projection"
	"frontend/components/texture"
	"frontend/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/materials/texturematerial"
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
	width, height := ioc.Get[window.Api](c).Window().GetSize()
	camera := world.NewEntity()
	world.SaveComponent(camera, transform.NewTransform())

	ortho := true
	ortho = false
	if ortho { // false for ortho
		projectionDepth := 2000
		world.SaveComponent(camera, projection.NewProjection(mgl32.Ortho(
			-float32(width)/2,
			float32(width)/2,
			-float32(height)/2,
			float32(height)/2,
			-float32(projectionDepth)/2,
			float32(projectionDepth)/2,
		)))
	} else {
		fovY := mgl32.DegToRad(90)
		aspectRatio := float32(width) / float32(height)
		var near, far float32 = 0.001, 1000
		world.SaveComponent(camera, projection.NewProjection(mgl32.Perspective(
			fovY,
			aspectRatio,
			near,
			far,
		)))
	}

	entity := world.NewEntity()

	world.SaveComponent(entity, transform.NewTransform().
		SetPos(transform.NewPos(0, 0, 300)).
		SetSize(transform.NewSize(100, 100, 100)))
	world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
	world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial))
	world.SaveComponent(entity, texture.NewTexture(TextureAssetID))
	world.SaveComponent(entity, ChangeTransformOverTimeComponent{})

	system := ChangeTransformOverTimeSystem{World: world}

	events.Listen(b, system.Update)

}

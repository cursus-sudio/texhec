package triangle

import (
	_ "embed"
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/materials/texturematerial"
	inputssystem "frontend/engine/systems/inputs"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MeshAssetID     assets.AssetID = "vao_asset"
	HalfMeshAssetID assets.AssetID = "half_vao_asset"
	TextureAssetID  assets.AssetID = "texture_asset"
)

func AddToWorld(c ioc.Dic, world ecs.World, b events.Builder) {
	{
		width, height := ioc.Get[window.Api](c).Window().GetSize()
		camera := world.NewEntity()
		world.SaveComponent(camera, transform.NewTransform().
			SetPos(mgl32.Vec3{0, 0, -100}).
			SetRotation(mgl32.QuatRotate(mgl32.DegToRad(180), mgl32.Vec3{1, 0, 0})))
		var orthoNear, orthoFar float32 = -1000, 1000
		world.SaveComponent(camera, projection.NewOrtho(float32(width), float32(height), orthoNear, orthoFar))

		fovY := mgl32.DegToRad(90)
		aspectRatio := float32(width) / float32(height)
		var perspectiveNear, perspectiveFar float32 = 0.001, 1000
		world.SaveComponent(camera, projection.NewPerspective(fovY, aspectRatio, perspectiveNear, perspectiveFar))

		events.Listen(b, func(e sdl.WindowEvent) {
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				width, height := e.Data1, e.Data2
				aspectRatio := float32(width) / float32(height)
				world.SaveComponent(camera, projection.NewPerspective(fovY, aspectRatio, perspectiveNear, perspectiveFar))
				world.SaveComponent(camera, projection.NewOrtho(float32(width), float32(height), orthoNear, orthoFar))
			}
		})
	}

	// {
	// 	entity := world.NewEntity()
	//
	// 	world.SaveComponent(entity, transform.NewTransform().
	// 		SetPos(mgl32.Vec3{0, 0, 300}).
	// 		SetSize(mgl32.Vec3{100, 100, 100}))
	// 	world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
	// 	// world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial3D, texturematerial.TextureMaterial2D))
	// 	// world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial2D))
	// 	world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial3D))
	// 	world.SaveComponent(entity, texture.NewTexture(TextureAssetID))
	// 	world.SaveComponent(entity, ChangeTransformOverTimeComponent{})
	// }

	system := ChangeTransformOverTimeSystem{World: world}
	events.Listen(b, system.Update)

	resizeSystem := inputssystem.NewResizeSystem()
	events.Listen(b, resizeSystem.Listen)

	{
		entity := world.NewEntity()

		world.SaveComponent(entity, transform.NewTransform().
			SetPos(mgl32.Vec3{0, 0, 100}).
			SetSize(mgl32.Vec3{1, 10000, 1}))
		world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
		world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial3D))
		// world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial2D))
		world.SaveComponent(entity, texture.NewTexture(TextureAssetID))

		onClickSystem := inputssystem.NewOnClickSystem[projection.Perspective](
			// onClickSystem := inputssystem.NewOnClickSystem[projection.Ortho](
			world,
			ioc.Get[window.Api](c),
			4,
			func() ecs.EntityId { return entity },
		)

		events.Listen(b, func(event sdl.MouseMotionEvent) {
			if err := onClickSystem.Listen(event); err != nil {
				panic(err)
			}
		})
	}
}

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
	"frontend/engine/systems/toysystems"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
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
			SetRotation(mgl32.QuatRotate(mgl32.DegToRad(180), mgl32.Vec3{1, 0, 0})),
		)
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

	// { // cube
	// 	entity := world.NewEntity()
	//
	// 	world.SaveComponent(entity, transform.NewTransform().
	// 		SetPos(mgl32.Vec3{0, 0, 300}).
	// 		SetSize(mgl32.Vec3{100, 100, 100}))
	// 	world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
	// 	world.SaveComponent(entity, material.NewMaterial(
	// 		texturematerial.TextureMaterial3D,
	// 		// texturematerial.TextureMaterial2D,
	// 	))
	// 	world.SaveComponent(entity, texture.NewTexture(TextureAssetID))
	// 	world.SaveComponent(entity, ChangeTransformOverTimeComponent{})
	// 	events.Listen(b, (&ChangeTransformOverTimeSystem{World: world}).Update)
	// }

	events.Listen(b, inputssystem.NewResizeSystem().Listen)

	{
		entity := world.NewEntity()

		world.SaveComponent(entity, transform.NewTransform().
			SetPos(mgl32.Vec3{0, 0, 100}).
			SetSize(mgl32.Vec3{20, 20, 10000}))
		world.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
		world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial3D))
		// world.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial2D))
		world.SaveComponent(entity, texture.NewTexture(TextureAssetID))

		onClickSystem := toysystems.NewShootRaySystem[projection.Perspective](
			// onClickSystem := inputssystem.NewOnClickSystem[projection.Ortho](
			world,
			ioc.Get[window.Api](c),
			ioc.Get[console.Console](c),
			100,
			func() ecs.EntityId { return entity },
		)

		events.Listen(b, func(event toysystems.ShootRayEvent) {
			if err := onClickSystem.Listen(event); err != nil {
				panic(err)
			}
		})
	}

	{
		events.Listen(b, func(event sdl.MouseButtonEvent) {
			events.Emit(
				b.Events(),
				toysystems.NewShootRayEvent(int(event.X), int(event.Y)),
			)
		})
	}

	{
		wPressed := false
		aPressed := false
		sPressed := false
		dPressed := false

		moveCameraSystem := func(event frames.FrameEvent) error {
			xAxis := 0
			if dPressed {
				xAxis = 1
			} else if aPressed {
				xAxis = -1
			}
			yAxis := 0
			if wPressed {
				yAxis = 1
			} else if sPressed {
				yAxis = -1
			}

			cameras := world.GetEntitiesWithComponents(ecs.GetComponentType(projection.Perspective{}))
			if len(cameras) != 1 {
				return projection.ErrWorldShouldHaveOneProjection
			}
			camera := cameras[0]
			var cameraTransform transform.Transform
			if err := world.GetComponent(camera, &cameraTransform); err != nil {
				return err
			}
			rotation := cameraTransform.Rotation
			mul := 100 * float32(event.Delta.Seconds())
			rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(xAxis)), mgl32.Vec3{0, 1, 0}))
			rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(yAxis)), mgl32.Vec3{-1, 0, 0}))
			cameraTransform.Rotation = rotation

			if err := world.SaveComponent(camera, cameraTransform); err != nil {
				return err
			}
			return nil
		}

		events.Listen(b, func(event frames.FrameEvent) {
			if err := moveCameraSystem(event); err != nil {
				panic(err)
			}

		})

		events.Listen(b, func(event sdl.KeyboardEvent) {
			pressed := event.State == sdl.PRESSED
			switch event.Keysym.Sym {
			case sdl.K_w:
				wPressed = pressed
				break
			case sdl.K_a:
				aPressed = pressed
				break
			case sdl.K_s:
				sPressed = pressed
				break
			case sdl.K_d:
				dPressed = pressed
				break
			}
		})
	}
}

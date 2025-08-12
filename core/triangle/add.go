package triangle

import (
	_ "embed"
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/materials/texturematerial"
	"frontend/engine/systems/projections"
	"frontend/services/assets"
	"frontend/services/colliders"
	"frontend/services/colliders/shapes"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/scenes"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MeshAssetID    assets.AssetID = "vao_asset"
	TextureAssetID assets.AssetID = "texture_asset"
)

func AddToWorld[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			camera := ctx.World.NewEntity()
			ctx.World.SaveComponent(camera, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, -100}).
				SetRotation(mgl32.QuatRotate(mgl32.DegToRad(180), mgl32.Vec3{1, 0, 0})),
			)
			ctx.World.SaveComponent(camera, projection.NewDynamicOrtho(
				-1000,
				+1000,
			))
			ctx.World.SaveComponent(camera, projection.NewDynamicPerspective(
				mgl32.DegToRad(90),
				0.01,
				1000,
			))
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projections.NewUpdateProjectionsEvent())
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) { // cube
			entity := ctx.World.NewEntity()
			ctx.World.SaveComponent(entity, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, 300}).
				SetSize(mgl32.Vec3{100, 100, 100}))
			ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
			ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
			// world.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
			ctx.World.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial))
			ctx.World.SaveComponent(entity, texture.NewTexture(TextureAssetID))
			ctx.World.SaveComponent(entity, ChangeTransformOverTimeComponent{})
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Listen(ctx.EventsBuilder, (&ChangeTransformOverTimeSystem{World: ctx.World}).Update)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			type OnHoveredDomainEvent struct{}
			type OnClickDomainEvent struct{}

			rectEntity := ctx.World.NewEntity()
			ctx.World.SaveComponent(rectEntity, transform.NewTransform().
				SetPos([3]float32{0, 0, 0}).
				SetSize([3]float32{100, 100, 0}))
			ctx.World.SaveComponent(rectEntity, mesh.NewMesh(MeshAssetID))
			ctx.World.SaveComponent(rectEntity, projection.NewUsedProjection[projection.Ortho]())
			// world.SaveComponent(rectEntity, projection.NewUsedProjection[projection.Perspective]())
			ctx.World.SaveComponent(rectEntity, material.NewMaterial(texturematerial.TextureMaterial))
			ctx.World.SaveComponent(rectEntity, texture.NewTexture(TextureAssetID))
			ctx.World.SaveComponent(rectEntity, mouse.NewMouseEvents().
				AddLeftClickEvents(OnClickDomainEvent{}),
			)
			ctx.World.SaveComponent(rectEntity, colliders.NewCollider([]colliders.Shape{
				shapes.NewRect2D(transform.NewTransform().SetSize([3]float32{1, 1}))}))

			events.Listen(ctx.EventsBuilder, func(e OnHoveredDomainEvent) {
				ioc.Get[console.Console](c).Print("damn it really is hovered\n")
			})

			events.Listen(ctx.EventsBuilder, func(fe frames.FrameEvent) {
				hovered, _ := ctx.World.GetComponentByType(rectEntity, ecs.GetComponentType(mouse.Hovered{}))
				if hovered == nil {
					return
				}
				ioc.Get[console.Console](c).Print("emit on frame should we not get one\n")
				events.Emit(ctx.Events, OnHoveredDomainEvent{})
			})

			events.Emit(ctx.Events, OnClickDomainEvent{})

			events.Listen(ctx.EventsBuilder, func(e OnClickDomainEvent) {
				ioc.Get[console.Console](c).PrintPermanent("damn it really is clicked\n")
			})
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			// move camera system inline
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

				cameras := ctx.World.GetEntitiesWithComponents(ecs.GetComponentType(projection.Perspective{}))
				if len(cameras) != 1 {
					return projection.ErrWorldShouldHaveOneProjection
				}
				camera := cameras[0]
				var cameraTransform transform.Transform
				if err := ctx.World.GetComponents(camera, &cameraTransform); err != nil {
					return err
				}
				rotation := cameraTransform.Rotation
				mul := 100 * float32(event.Delta.Seconds())
				rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(xAxis)), mgl32.Vec3{0, 1, 0}))
				rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(yAxis)), mgl32.Vec3{-1, 0, 0}))
				cameraTransform.Rotation = rotation

				if err := ctx.World.SaveComponent(camera, cameraTransform); err != nil {
					return err
				}
				return nil
			}

			events.ListenE(ctx.EventsBuilder, moveCameraSystem)

			events.Listen(ctx.EventsBuilder, func(event sdl.KeyboardEvent) {
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
		})
		return b
	})
}

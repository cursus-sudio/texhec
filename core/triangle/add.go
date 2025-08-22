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
				SetPos(mgl32.Vec3{0, 0, 100}).
				SetRotation(mgl32.QuatIdent()),
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
		b.OnLoad(func(ctx scenes.SceneCtx) {
			entity := ctx.World.NewEntity()
			ctx.World.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial))
			ctx.World.SaveComponent(entity, texturematerial.NewWorldTextureMaterialComponent(
				[]assets.AssetID{TextureAssetID},
				[]assets.AssetID{MeshAssetID},
			))
		})

		b.OnLoad(func(ctx scenes.SceneCtx) { // cube
			entity := ctx.World.NewEntity()
			ctx.World.SaveComponent(entity, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, -300}).
				SetSize(mgl32.Vec3{100, 100, 100}))
			ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
			ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
			// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
			ctx.World.SaveComponent(entity, texturematerial.TextureMaterialComponent{})
			ctx.World.SaveComponent(entity, texture.NewTexture(TextureAssetID))
			ctx.World.SaveComponent(entity, ChangeTransformOverTimeComponent{})
		})
		b.OnLoad(func(ctx scenes.SceneCtx) {
			projections.NewOcclusionSystem(ctx.World)
			system := NewChangeTransformOverTimeSystem(ctx.World)
			events.Listen(ctx.EventsBuilder, system.Update)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			type Marked struct{}

			type OnHoveredDomainEvent struct{}
			type OnClickDomainEvent struct{}

			{
				events.Listen(ctx.EventsBuilder, func(e OnHoveredDomainEvent) {
					ioc.Get[console.Console](c).Print("damn it really is hovered\n")
				})

				events.Listen(ctx.EventsBuilder, func(e OnClickDomainEvent) {
					ioc.Get[console.Console](c).PrintPermanent("damn it really is clicked\n")
				})

				liveQuery := ctx.World.QueryEntitiesWithComponents(
					ecs.GetComponentType(Marked{}),
					ecs.GetComponentType(mouse.Hovered{}),
				)
				events.Listen(ctx.EventsBuilder, func(fe frames.FrameEvent) {
					for range liveQuery.Entities() {
						events.Emit(ctx.Events, OnHoveredDomainEvent{})
					}
				})
			}

			rows := 100
			cols := 1000
			for i := 0; i < rows*cols; i++ {
				row := i / cols
				col := i % cols
				entity := ctx.World.NewEntity()
				ctx.World.SaveComponent(entity, transform.NewTransform().
					SetPos([3]float32{float32(col) * 101, float32(row) * 101, 0}).
					SetSize([3]float32{100, 100, 1}))
				ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
				ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
				// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
				// ctx.World.SaveComponent(entity, material.NewMaterial(texturematerial.TextureMaterial))
				ctx.World.SaveComponent(entity, texturematerial.TextureMaterialComponent{})
				ctx.World.SaveComponent(entity, texture.NewTexture(TextureAssetID))
				ctx.World.SaveComponent(entity, Marked{})
				ctx.World.SaveComponent(entity, mouse.NewMouseEvents().AddLeftClickEvents(OnClickDomainEvent{}))
				ctx.World.SaveComponent(entity, colliders.NewCollider([]colliders.Shape{
					shapes.NewRect2D(transform.NewTransform().SetSize([3]float32{1, 1}))}))
			}
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			// move camera system inline
			wPressed := false
			aPressed := false
			sPressed := false
			dPressed := false
			camerasQuery := ctx.World.QueryEntitiesWithComponents(
				ecs.GetComponentType(projection.Perspective{}),
			)

			moveCameraSystem := func(event frames.FrameEvent) error {
				xAxis := 0
				if dPressed {
					xAxis = -1
				} else if aPressed {
					xAxis = 1
				}
				yAxis := 0
				if wPressed {
					yAxis = -1
				} else if sPressed {
					yAxis = 1
				}

				cameras := camerasQuery.Entities()
				if len(cameras) != 1 {
					return projection.ErrWorldShouldHaveOneProjection
				}
				camera := cameras[0]
				cameraTransform, err := ecs.GetComponent[transform.Transform](ctx.World, camera)
				if err != nil {
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

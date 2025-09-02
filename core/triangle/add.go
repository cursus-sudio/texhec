package triangle

import (
	_ "embed"
	"fmt"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/mainpipeline"
	"frontend/engine/systems/projections"
	"frontend/engine/tools/worldmesh"
	"frontend/engine/tools/worldprojections"
	"frontend/engine/tools/worldtexture"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MeshAssetID     assets.AssetID = "vao_asset"
	TextureAssetID  assets.AssetID = "texture_asset"
	ColliderAssetID assets.AssetID = "collider_asset"
)

func AddToWorld[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			camera := ctx.World.NewEntity()
			ctx.World.SaveComponent(camera, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, 100}),
			)
			ctx.World.SaveComponent(camera, projection.NewDynamicOrtho(
				-1000,
				+1000,
				1,
			))
			ctx.World.SaveComponent(camera, projection.NewDynamicPerspective(
				mgl32.DegToRad(90),
				0.01,
				1000,
			))
			// uiCamera := ctx.World.NewEntity()
			// ctx.World.SaveComponent(uiCamera, transform.NewTransform().
			// 	SetPos(mgl32.Vec3{0, -50, 100}),
			// )
			// ctx.World.SaveComponent(uiCamera, projection.NewDynamicOrtho(
			// 	-1000,
			// 	+1000,
			// 	1,
			// ))
			// ctx.World.SaveComponent(uiCamera, projection.NewDynamicPerspective(
			// 	mgl32.DegToRad(90),
			// 	0.01,
			// 	1000,
			// ))
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
		worldMeshFactory := ioc.Get[worldmesh.RegisterFactory[mainpipeline.Vertex]](c)
		worldTextureFactory := ioc.Get[worldtexture.RegisterFactory](c)
		b.OnLoad(func(ctx scenes.SceneCtx) {
			projectionsRegister := worldprojections.NewWorldProjectionsRegister(
				ecs.GetComponentType(projection.Ortho{}),
				ecs.GetComponentType(projection.Perspective{}),
			)
			ctx.World.SaveRegister(projectionsRegister)

			textureRegister, err := worldTextureFactory.New(TextureAssetID)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
			ctx.World.SaveRegister(textureRegister)

			meshRegister, err := worldMeshFactory.New(MeshAssetID)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
			ctx.World.SaveRegister(meshRegister)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) { // cube
			entity := ctx.World.NewEntity()
			ctx.World.SaveComponent(entity, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, -300}).
				SetSize(mgl32.Vec3{100, 100, 100}))
			ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
			ctx.World.SaveComponent(entity, texture.NewTexture(TextureAssetID))
			ctx.World.SaveComponent(entity, mainpipeline.PipelineComponent{})
			ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
			// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
			ctx.World.SaveComponent(entity, ChangeTransformOverTimeComponent{})
		})
		b.OnLoad(func(ctx scenes.SceneCtx) {
			pipeline, err := mainpipeline.NewSystem(
				ctx.World,
				ioc.Get[window.Api](c),
				ioc.Get[assets.AssetsStorage](c),
				ioc.Get[logger.Logger](c),
				[]ecs.ComponentType{
					ecs.GetComponentType(projection.Visible{}),
				},
			)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
			events.ListenE(ctx.EventsBuilder, pipeline.Listen)
			projections.NewOcclusionSystem(ctx.World)
			system := NewChangeTransformOverTimeSystem(ctx.World)
			events.Listen(ctx.EventsBuilder, system.Update)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) {
			type OnHoveredDomainEvent struct {
				entity   ecs.EntityID
				row, col int
			}
			type OnClickDomainEvent struct {
				entity   ecs.EntityID
				row, col int
			}

			{
				events.Listen(ctx.EventsBuilder, func(e OnHoveredDomainEvent) {
					ioc.Get[console.Console](c).Print(
						fmt.Sprintf("damn it really is hovered %v (%d, %d)\n", e.entity, e.col, e.row),
					)
				})

				events.Listen(ctx.EventsBuilder, func(e OnClickDomainEvent) {
					ioc.Get[console.Console](c).PrintPermanent(
						fmt.Sprintf("damn it really is clicked %v (%d, %d)\n", e.entity, e.col, e.row),
					)
				})
			}

			rows := 100
			cols := 100
			var size float32 = 100
			var gap float32 = 0
			for i := 0; i < rows*cols; i++ {
				row := i / cols
				col := i % cols
				entity := ctx.World.NewEntity()
				ctx.World.SaveComponent(entity, transform.NewTransform().
					SetPos([3]float32{float32(col) * (size + gap), float32(row) * (size + gap), 0}).
					SetSize([3]float32{size, size, 1}))
				ctx.World.SaveComponent(entity, transform.NewStatic())
				ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
				ctx.World.SaveComponent(entity, texture.NewTexture(TextureAssetID))
				ctx.World.SaveComponent(entity, mainpipeline.PipelineComponent{})
				ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
				// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
				ctx.World.SaveComponent(entity, collider.NewCollider(ColliderAssetID))
				ctx.World.SaveComponent(entity, mouse.NewMouseEvents().
					AddLeftClickEvents(OnClickDomainEvent{entity, row, col}).
					AddMouseHoverEvents(OnHoveredDomainEvent{entity, row, col}),
				)
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

				cameras := camerasQuery.Entities()
				for _, camera := range cameras {
					cameraTransform, err := ecs.GetComponent[transform.Transform](ctx.World, camera)
					if err != nil {
						return err
					}
					{
						pos := cameraTransform.Pos
						mul := 1000 * float32(event.Delta.Seconds())
						pos[0] += mul * float32(xAxis)
						pos[1] += mul * float32(yAxis)
						cameraTransform.Pos = pos
					}
					// rotation := cameraTransform.Rotation
					// mul := 100 * float32(event.Delta.Seconds())
					// rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(xAxis)), mgl32.Vec3{0, 1, 0}))
					// rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(yAxis)), mgl32.Vec3{-1, 0, 0}))
					// cameraTransform.Rotation = rotation

					if err := ctx.World.SaveComponent(camera, cameraTransform); err != nil {
						return err
					}
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

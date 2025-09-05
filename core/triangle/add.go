package triangle

import (
	_ "embed"
	"fmt"
	"frontend/engine/components/anchor"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/anchorsystem"
	"frontend/engine/systems/mainpipeline"
	"frontend/engine/systems/projections"
	"frontend/engine/systems/transformsystem"
	"frontend/engine/tools/worldmesh"
	"frontend/engine/tools/worldprojections"
	"frontend/engine/tools/worldtexture"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"math"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func AddToWorld[SceneBuilder scenes.SceneBuilder](b ioc.Builder) {
	type MobileCamera struct{}
	ioc.WrapService(b, scenes.LoadWorld, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			camera := ctx.World.NewEntity()
			ctx.World.SaveComponent(camera, transform.NewTransform())
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
			ctx.World.SaveComponent(camera, MobileCamera{})
			uiCamera := ctx.World.NewEntity()
			ctx.World.SaveComponent(uiCamera, transform.NewTransform().
				SetPos(mgl32.Vec3{0, 0, 10000}))
			ctx.World.SaveComponent(uiCamera, projection.NewDynamicOrtho(
				-1,
				+1,
				1,
			))

			type QuitEvent struct{}

			events.Listen(ctx.EventsBuilder, func(e QuitEvent) {
				ioc.Get[runtime.Runtime](c).Stop()
			})
			exitBtn := ctx.World.NewEntity()
			ctx.World.SaveComponent(exitBtn, transform.NewTransform().
				SetSize(mgl32.Vec3{100, 100, 1}))
			ctx.World.SaveComponent(exitBtn, anchor.NewParentAnchor(uiCamera).
				SetPivotPoint(mgl32.Vec3{0, 1, .5}))
			ctx.World.SaveComponent(exitBtn, transform.NewPivotPoint(mgl32.Vec3{1, 0, .5}))
			ctx.World.SaveComponent(exitBtn, mesh.NewMesh(MeshAssetID))
			ctx.World.SaveComponent(exitBtn, texture.NewTexture(Texture4AssetID))
			ctx.World.SaveComponent(exitBtn, mainpipeline.PipelineComponent{})
			ctx.World.SaveComponent(exitBtn, projection.NewUsedProjection[projection.Ortho]())
			ctx.World.SaveComponent(exitBtn, collider.NewCollider(ColliderAssetID))
			ctx.World.SaveComponent(exitBtn, mouse.NewMouseEvents().
				// AddMouseHoverEvents(QuitEvent{}).
				AddLeftClickEvents(QuitEvent{}),
			)
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadInitialEvents, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			events.Emit(ctx.Events, projections.NewUpdateProjectionsEvent())
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadFirst, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			anchorsystem.NewAnchorSystem(ctx.World)
			transformsystem.NewPivotPointSystem(ctx.World)
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

			textureRegister, err := worldTextureFactory.New(
				Texture1AssetID,
				Texture2AssetID,
				Texture3AssetID,
				Texture4AssetID,
			)
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
			ctx.World.SaveComponent(entity, texture.NewTexture(Texture2AssetID))
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
				[]ecs.ComponentType{},
			)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
			events.ListenE(ctx.EventsBuilder, pipeline.Listen)
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
					SetPos([3]float32{float32(col) * (size + gap), float32(row) * (size + gap), 10}).
					SetSize([3]float32{size, size, 1}))
				ctx.World.SaveComponent(entity, transform.NewStatic())
				ctx.World.SaveComponent(entity, mesh.NewMesh(MeshAssetID))
				ctx.World.SaveComponent(entity, texture.NewTexture(Texture1AssetID))
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
				ecs.GetComponentType(transform.Transform{}),
				ecs.GetComponentType(projection.DynamicOrtho{}),
				ecs.GetComponentType(MobileCamera{}),
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

			events.ListenE(ctx.EventsBuilder, func(event sdl.MouseWheelEvent) error {
				if event.Y == 0 {
					return nil
				}
				cameras := camerasQuery.Entities()
				var mul = float32(math.Pow(10, float64(event.Y)/50))
				for _, camera := range cameras {
					ortho, err := ecs.GetComponent[projection.DynamicOrtho](ctx.World, camera)
					if err != nil {
						return err
					}

					ortho.Zoom *= mul
					ortho.Zoom = max(min(ortho.Zoom, 5), 0.1)

					if err := ctx.World.SaveComponent(camera, ortho); err != nil {
						return err
					}
				}
				return nil
			})

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

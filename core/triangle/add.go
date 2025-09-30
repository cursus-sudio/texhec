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
	"frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/projections"
	"frontend/engine/systems/transformsystem"
	"frontend/engine/tools/worldprojections"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/frames"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"frontend/services/scenes"
	"math"
	"shared/services/ecs"
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
			ecs.SaveComponent(ctx.World.Components(), camera, transform.NewTransform())
			ecs.SaveComponent(ctx.World.Components(), camera, projection.NewDynamicOrtho(
				-1000,
				+1000,
				1,
			))
			ecs.SaveComponent(ctx.World.Components(), camera, projection.NewDynamicPerspective(
				mgl32.DegToRad(90),
				0.01,
				1000,
			))
			ecs.SaveComponent(ctx.World.Components(), camera, MobileCamera{})
			uiCamera := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), uiCamera, transform.NewTransform().Ptr().
				SetPos(mgl32.Vec3{0, 0, 10000}).Val())
			ecs.SaveComponent(ctx.World.Components(), uiCamera, projection.NewDynamicOrtho(
				-1000,
				+1000,
				1,
			))

			type QuitEvent struct{}

			events.Listen(ctx.EventsBuilder, func(e QuitEvent) {
				ioc.Get[runtime.Runtime](c).Stop()
			})
			exitBtn := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), exitBtn, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 100, 1}).Val())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, anchor.NewParentAnchor(uiCamera).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 1, .5}).
				Val())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, transform.NewPivotPoint(mgl32.Vec3{1, 0, .5}))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, mesh.NewMesh(MeshAssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, texture.NewTexture(Texture4AssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, genericrenderer.PipelineComponent{})
			ecs.SaveComponent(ctx.World.Components(), exitBtn, projection.NewUsedProjection[projection.Ortho]())
			ecs.SaveComponent(ctx.World.Components(), exitBtn, collider.NewCollider(ColliderAssetID))
			ecs.SaveComponent(ctx.World.Components(), exitBtn, mouse.NewMouseEvents().
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
			anchorsystem.NewAnchorSystem(ctx.World, ioc.Get[logger.Logger](c))
			transformsystem.NewPivotPointSystem(ctx.World, ioc.Get[logger.Logger](c))
		})
		return b
	})

	ioc.WrapService(b, scenes.LoadDomain, func(c ioc.Dic, b SceneBuilder) SceneBuilder {
		b.OnLoad(func(ctx scenes.SceneCtx) {
			projectionsRegister := worldprojections.NewWorldProjectionsRegister(
				ecs.GetComponentType(projection.Ortho{}),
				ecs.GetComponentType(projection.Perspective{}),
			)
			ctx.World.SaveRegister(projectionsRegister)
		})

		b.OnLoad(func(ctx scenes.SceneCtx) { // cube
			entity := ctx.World.NewEntity()
			ecs.SaveComponent(ctx.World.Components(), entity, transform.NewTransform().Ptr().
				SetPos(mgl32.Vec3{0, 0, -300}).
				SetSize(mgl32.Vec3{100, 100, 100}).Val())
			ecs.SaveComponent(ctx.World.Components(), entity, mesh.NewMesh(MeshAssetID))
			ecs.SaveComponent(ctx.World.Components(), entity, texture.NewTexture(Texture2AssetID))
			ecs.SaveComponent(ctx.World.Components(), entity, genericrenderer.PipelineComponent{})
			ecs.SaveComponent(ctx.World.Components(), entity, projection.NewUsedProjection[projection.Perspective]())
			// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Ortho]())
			ecs.SaveComponent(ctx.World.Components(), entity, ChangeTransformOverTimeComponent{})
		})
		b.OnLoad(func(ctx scenes.SceneCtx) {
			pipeline, err := genericrenderer.NewSystem(
				ctx.World,
				ioc.Get[window.Api](c),
				ioc.Get[assets.AssetsStorage](c),
				ioc.Get[logger.Logger](c),
				ioc.Get[vbo.VBOFactory[genericrenderer.Vertex]](c),
				[]ecs.ComponentType{},
			)
			if err != nil {
				ioc.Get[logger.Logger](c).Error(err)
			}
			events.ListenE(ctx.EventsBuilder, pipeline.Listen)
			system := NewChangeTransformOverTimeSystem(ctx.World, ioc.Get[logger.Logger](c))
			events.Listen(ctx.EventsBuilder, system.Listen)
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

			rows := 10
			cols := 10
			var size float32 = 100
			var gap float32 = 0
			for i := 0; i < rows*cols; i++ {
				row := i / cols
				col := i % cols
				entity := ctx.World.NewEntity()
				ecs.SaveComponent(ctx.World.Components(), entity, transform.NewTransform().Ptr().
					SetPos([3]float32{float32(col) * (size + gap), float32(row) * (size + gap), 0}).
					SetSize([3]float32{size, size, 1}).Val())
				ecs.SaveComponent(ctx.World.Components(), entity, transform.NewStatic())
				ecs.SaveComponent(ctx.World.Components(), entity, mesh.NewMesh(MeshAssetID))
				ecs.SaveComponent(ctx.World.Components(), entity, texture.NewTexture(Texture1AssetID))
				ecs.SaveComponent(ctx.World.Components(), entity, genericrenderer.PipelineComponent{})
				ecs.SaveComponent(ctx.World.Components(), entity, projection.NewUsedProjection[projection.Ortho]())
				// ctx.World.SaveComponent(entity, projection.NewUsedProjection[projection.Perspective]())
				ecs.SaveComponent(ctx.World.Components(), entity, collider.NewCollider(ColliderAssetID))
				ecs.SaveComponent(ctx.World.Components(), entity, mouse.NewMouseEvents().
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
			transformArray := ecs.GetComponentsArray[transform.Transform](ctx.World.Components())

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
					cameraTransform, err := transformArray.GetComponent(camera)
					if err != nil {
						return err
					}
					{
						pos := cameraTransform.Pos
						mul := 1000 * float32(event.Delta.Seconds())
						pos[0] += mul * float32(xAxis)
						pos[1] += mul * float32(yAxis)
						cameraTransform.SetPos(pos)
					}
					// rotation := cameraTransform.Rotation
					// mul := 100 * float32(event.Delta.Seconds())
					// rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(xAxis)), mgl32.Vec3{0, 1, 0}))
					// rotation = rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(mul*float32(yAxis)), mgl32.Vec3{-1, 0, 0}))
					// cameraTransform.Rotation = rotation

					if err := transformArray.SaveComponent(camera, cameraTransform); err != nil {
						return err
					}
				}
				return nil
			}

			events.ListenE(ctx.EventsBuilder, moveCameraSystem)

			dynamicOrthoArray := ecs.GetComponentsArray[projection.DynamicOrtho](ctx.World.Components())
			events.ListenE(ctx.EventsBuilder, func(event sdl.MouseWheelEvent) error {
				if event.Y == 0 {
					return nil
				}
				cameras := camerasQuery.Entities()
				var mul = float32(math.Pow(10, float64(event.Y)/50))
				for _, camera := range cameras {
					ortho, err := dynamicOrthoArray.GetComponent(camera)
					if err != nil {
						return err
					}

					ortho.Zoom *= mul
					ortho.Zoom = max(min(ortho.Zoom, 5), 0.1)

					if err := dynamicOrthoArray.SaveComponent(camera, ortho); err != nil {
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

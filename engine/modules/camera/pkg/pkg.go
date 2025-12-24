package camerapkg

import (
	"engine/modules/camera"
	"engine/modules/camera/internal/cameralimitsys"
	"engine/modules/camera/internal/cameratool"
	"engine/modules/camera/internal/mobilecamerasys"
	"engine/modules/camera/internal/projectionsys"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct {
	minZoom, maxZoom float32
}

func Package(minZoom, maxZoom float32) ioc.Pkg {
	return pkg{
		minZoom: minZoom,
		maxZoom: maxZoom,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// camera components
			Register(camera.Component{}).
			Register(camera.MobileCameraComponent{}).
			Register(camera.CameraLimitsComponent{}).
			Register(camera.ViewportComponent{}).
			Register(camera.NormalizedViewportComponent{}).
			// projections components
			Register(camera.OrthoComponent{}).
			Register(camera.OrthoResolutionComponent{}).
			Register(camera.PerspectiveComponent{}).
			Register(camera.DynamicPerspectiveComponent{}).
			// events
			Register(camera.ChangedResolutionEvent{})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) cameratool.CameraResolverFactory {
		return cameratool.NewCameraResolverFactory()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.ToolFactory {
		return ioc.Get[cameratool.CameraResolverFactory](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.CameraUp { return camera.CameraUp(mgl32.Vec3{0, 1, 0}) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.CameraForward { return camera.CameraForward(mgl32.Vec3{0, 0, -1}) })

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s cameratool.CameraResolverFactory) cameratool.CameraResolverFactory {
		viewport := func(w ecs.World) func(entity ecs.EntityID) func() (x, y, w, h int32) {
			window := ioc.Get[window.Api](c)
			viewportArray := ecs.GetComponentsArray[camera.ViewportComponent](w)
			normalizedViewportArray := ecs.GetComponentsArray[camera.NormalizedViewportComponent](w)
			return func(entity ecs.EntityID) func() (x int32, y int32, w int32, h int32) {
				return func() (rx int32, ry int32, rw int32, rh int32) { // r from result
					viewportComponent, ok := viewportArray.Get(entity)
					if ok {
						return viewportComponent.Viewport()
					}
					normalizedViewportComponent, ok := normalizedViewportArray.Get(entity)
					if ok {
						return normalizedViewportComponent.Viewport(window.Window().GetSize())
					}

					w, h := window.Window().GetSize()
					return 0, 0, w, h
				}
			}
		}
		s.Register(reflect.TypeFor[camera.OrthoComponent](), func(world camera.World) func(entity ecs.EntityID) (camera.Object, error) {
			orthoArray := ecs.GetComponentsArray[camera.OrthoComponent](world)
			orthoResolutionArray := ecs.GetComponentsArray[camera.OrthoResolutionComponent](world)
			viewport := viewport(world)
			return func(entity ecs.EntityID) (camera.Object, error) {
				viewport := viewport(entity)
				getCameraTransformMatrix := func() mgl32.Mat4 {
					pos, _ := world.Transform().AbsolutePos().Get(entity)
					rot, _ := world.Transform().AbsoluteRotation().Get(entity)

					cameraRot := rot.Rotation.Inverse()
					cameraPos := rot.Rotation.Rotate(pos.Pos.Mul(-1))
					return cameraRot.Mat4().Mul4(mgl32.Translate3D(cameraPos.X(), cameraPos.Y(), cameraPos.Z()))
				}
				getProjection := func() camera.OrthoComponent {
					ortho, ok := orthoArray.Get(entity)
					if !ok {
						return ortho
					}
					return ortho
				}
				getProjectionMatrix := func() mgl32.Mat4 {
					p := getProjection()
					orthoResolution, ok := orthoResolutionArray.Get(entity)
					if !ok {
						orthoResolution = camera.GetViewportOrthoResolution(viewport())
					}
					return p.GetMatrix(orthoResolution.Elem())
				}

				cameraGroups, ok := ecs.GetComponent[groups.GroupsComponent](world, entity)
				if !ok {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameratool.NewObject(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					viewport,
					func(mousePos window.MousePos) collider.Ray {
						return mobilecamerasys.ShootRay(
							getProjectionMatrix(),
							getCameraTransformMatrix(),
							mousePos,
							viewport,
							nil,
						)
					},
					cameraGroups,
				)
				return camera, nil
			}
		})

		//

		s.Register(reflect.TypeFor[camera.PerspectiveComponent](), func(world camera.World) func(entity ecs.EntityID) (camera.Object, error) {
			perspectiveArray := ecs.GetComponentsArray[camera.PerspectiveComponent](world)
			viewport := viewport(world)
			return func(entity ecs.EntityID) (camera.Object, error) {
				viewport := viewport(entity)
				getCameraTransformMatrix := func() mgl32.Mat4 {
					pos, _ := world.Transform().AbsolutePos().Get(entity)
					rot, _ := world.Transform().AbsoluteRotation().Get(entity)

					up, forward := ioc.Get[camera.CameraUp](c), ioc.Get[camera.CameraForward](c)
					return mgl32.LookAtV(
						pos.Pos,
						pos.Pos.Add(rot.Rotation.Rotate(mgl32.Vec3(forward))),
						mgl32.Vec3(up),
					)
				}
				getProjection := func() camera.PerspectiveComponent {
					perspective, ok := perspectiveArray.Get(entity)
					if !ok {
						return perspective
					}
					return perspective
				}
				getProjectionMatrix := func() mgl32.Mat4 {
					p := getProjection()
					return mgl32.Perspective(p.FovY, p.AspectRatio, p.Near, p.Far)
				}

				cameraGroups, ok := ecs.GetComponent[groups.GroupsComponent](world, entity)
				if !ok {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameratool.NewObject(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					viewport,
					func(mousePos window.MousePos) collider.Ray {
						pos, _ := world.Transform().AbsolutePos().Get(entity)
						return mobilecamerasys.ShootRay(
							getProjectionMatrix(),
							getCameraTransformMatrix(),
							mousePos,
							viewport,
							&pos.Pos,
						)
					},
					cameraGroups,
				)
				return camera, nil
			}
		})

		return s
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.System {
		return ecs.NewSystemRegister(func(w camera.World) error {
			logger := ioc.Get[logger.Logger](c)
			ecs.RegisterSystems(w,
				ecs.NewSystemRegister(func(w ecs.World) error {
					cameraArray := ecs.GetComponentsArray[camera.Component](w)
					orthoArray := ecs.GetComponentsArray[camera.OrthoComponent](w)
					perspectiveArray := ecs.GetComponentsArray[camera.PerspectiveComponent](w)

					cameraArray.AddDependency(orthoArray)
					cameraArray.AddDependency(perspectiveArray)

					orthoDirtySet := ecs.NewDirtySet()
					orthoArray.AddDirtySet(orthoDirtySet)

					cameraArray.BeforeGet(func() {
						entities := orthoDirtySet.Get()
						for _, entity := range entities {
							cameraArray.Set(entity, camera.NewCamera[camera.OrthoComponent]())
						}
					})

					perspectiveDirtySet := ecs.NewDirtySet()
					perspectiveArray.AddDirtySet(perspectiveDirtySet)

					cameraArray.BeforeGet(func() {
						entities := perspectiveDirtySet.Get()
						for _, entity := range entities {
							cameraArray.Set(entity, camera.NewCamera[camera.PerspectiveComponent]())
						}
					})

					events.Listen(w.EventsBuilder(), func(e sdl.WindowEvent) {
						if e.Event == sdl.WINDOWEVENT_RESIZED {
							events.Emit(w.Events(), camera.NewUpdateProjectionsEvent())
						}
					})
					return nil
				}),
				// todo change this to change ortho and size according to viewport
				projectionsys.NewUpdateProjectionsSystem(
					ioc.Get[window.Api](c),
					logger,
					ioc.Get[camera.ToolFactory](c),
				),
				mobilecamerasys.NewScrollSystem(
					logger,
					ioc.Get[camera.ToolFactory](c),
					ioc.Get[window.Api](c),
					pkg.minZoom, pkg.maxZoom, // min and max zoom
				),
				mobilecamerasys.NewDragSystem(
					sdl.BUTTON_LEFT,
					ioc.Get[camera.ToolFactory](c),
					ioc.Get[window.Api](c),
					logger,
				),
				mobilecamerasys.NewWasdSystem(
					logger,
					ioc.Get[camera.ToolFactory](c),
					1.0, // speed
				),
				cameralimitsys.NewOrthoSys(
					ioc.Get[camera.ToolFactory](c),
					logger,
				),
			)
			return nil
		})
	})
}

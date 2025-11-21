package camerapkg

import (
	"frontend/modules/camera"
	"frontend/modules/camera/internal/cameralimitsys"
	"frontend/modules/camera/internal/cameratool"
	"frontend/modules/camera/internal/mobilecamerasys"
	"frontend/modules/camera/internal/projectionsys"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"frontend/modules/transform"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

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
	ioc.RegisterSingleton(b, func(c ioc.Dic) cameratool.CameraResolverFactory {
		return cameratool.NewCameraResolverFactory()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[camera.CameraTool] {
		return ioc.Get[cameratool.CameraResolverFactory](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.CameraUp { return camera.CameraUp(mgl32.Vec3{0, 1, 0}) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) camera.CameraForward { return camera.CameraForward(mgl32.Vec3{0, 0, -1}) })

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s cameratool.CameraResolverFactory) cameratool.CameraResolverFactory {
		s.Register(ecs.GetComponentType(camera.OrthoComponent{}), func(world ecs.World) func(entity ecs.EntityID) (camera.CameraService, error) {
			logger := ioc.Get[logger.Logger](c)
			transformTransaction := ioc.Get[ecs.ToolFactory[transform.TransformTool]](c).Build(world).Transaction()
			orthoArray := ecs.GetComponentsArray[camera.OrthoComponent](world)
			return func(entity ecs.EntityID) (camera.CameraService, error) {
				transform := transformTransaction.GetEntity(entity)
				getCameraTransformMatrix := func() mgl32.Mat4 {
					pos, err := transform.AbsolutePos().Get()
					if err != nil {
						logger.Warn(err)
						return mgl32.Mat4{}
					}
					rot, err := transform.AbsoluteRotation().Get()
					if err != nil {
						logger.Warn(err)
						return mgl32.Mat4{}
					}

					cameraRot := rot.Rotation.Inverse()
					cameraPos := rot.Rotation.Rotate(pos.Pos.Mul(-1))
					return cameraRot.Mat4().Mul4(mgl32.Translate3D(cameraPos.X(), cameraPos.Y(), cameraPos.Z()))
				}
				getProjection := func() camera.OrthoComponent {
					ortho, err := orthoArray.GetComponent(entity)
					if err != nil {
						return ortho
					}
					return ortho
				}
				getProjectionMatrix := func() mgl32.Mat4 {
					p := getProjection()
					return mgl32.Ortho(
						-p.Width/2, p.Width/2,
						-p.Height/2, p.Height/2,
						p.Near, p.Far,
					)
				}

				cameraGroups, err := ecs.GetComponent[groups.GroupsComponent](world, entity)
				if err != nil {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameratool.NewCameraService(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					func(mousePos mgl32.Vec2) collider.Ray {
						return mobilecamerasys.ShootRay(
							getProjectionMatrix(),
							getCameraTransformMatrix(),
							mousePos,
							nil,
						)
					},
					cameraGroups,
				)
				return camera, nil
			}
		})

		//

		s.Register(ecs.GetComponentType(camera.Perspective{}), func(world ecs.World) func(entity ecs.EntityID) (camera.CameraService, error) {
			logger := ioc.Get[logger.Logger](c)
			transformTransaction := ioc.Get[ecs.ToolFactory[transform.TransformTool]](c).Build(world).Transaction()
			perspectiveArray := ecs.GetComponentsArray[camera.Perspective](world)

			return func(entity ecs.EntityID) (camera.CameraService, error) {
				transform := transformTransaction.GetEntity(entity)
				getCameraTransformMatrix := func() mgl32.Mat4 {
					pos, err := transform.AbsolutePos().Get()
					if err != nil {
						logger.Warn(err)
						return mgl32.Mat4{}
					}
					rot, err := transform.AbsoluteRotation().Get()
					if err != nil {
						logger.Warn(err)
						return mgl32.Mat4{}
					}

					up, forward := ioc.Get[camera.CameraUp](c), ioc.Get[camera.CameraForward](c)
					return mgl32.LookAtV(
						pos.Pos,
						pos.Pos.Add(rot.Rotation.Rotate(mgl32.Vec3(forward))),
						mgl32.Vec3(up),
					)
				}
				getProjection := func() camera.Perspective {
					perspective, err := perspectiveArray.GetComponent(entity)
					if err != nil {
						return perspective
					}
					return perspective
				}
				getProjectionMatrix := func() mgl32.Mat4 {
					p := getProjection()
					return mgl32.Perspective(p.FovY, p.AspectRatio, p.Near, p.Far)
				}

				cameraGroups, err := ecs.GetComponent[groups.GroupsComponent](world, entity)
				if err != nil {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameratool.NewCameraService(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					func(mousePos mgl32.Vec2) collider.Ray {
						pos, err := transform.AbsolutePos().Get()
						if err != nil {
							logger.Warn(err)
							return collider.Ray{}
						}
						return mobilecamerasys.ShootRay(
							getProjectionMatrix(),
							getCameraTransformMatrix(),
							mousePos,
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
		return ecs.NewSystemRegister(func(w ecs.World) error {
			logger := ioc.Get[logger.Logger](c)
			ecs.RegisterSystems(w,
				projectionsys.NewUpdateProjectionsSystem(
					ioc.Get[window.Api](c),
					logger,
					ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
				),
				mobilecamerasys.NewScrollSystem(
					logger,
					ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
					ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
					ioc.Get[window.Api](c),
					pkg.minZoom, pkg.maxZoom, // min and max zoom
				),
				mobilecamerasys.NewDragSystem(
					sdl.BUTTON_LEFT,
					ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
					ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
					ioc.Get[window.Api](c),
					logger,
				),
				mobilecamerasys.NewWasdSystem(
					logger,
					ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
					ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
					1.0, // speed
				),
				cameralimitsys.NewOrthoSys(
					ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
					logger,
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					cameraArray := ecs.GetComponentsArray[camera.CameraComponent](w)

					orthoArray := ecs.GetComponentsArray[camera.OrthoComponent](w)
					orthoArray.OnAdd(func(ei []ecs.EntityID) {
						t := cameraArray.Transaction()
						for _, e := range ei {
							t.SaveComponent(e, camera.NewCamera(ecs.GetComponentType(camera.OrthoComponent{})))
						}
						logger.Warn(ecs.FlushMany(t))
					})

					perspectiveArray := ecs.GetComponentsArray[camera.Perspective](w)
					perspectiveArray.OnAdd(func(ei []ecs.EntityID) {
						t := cameraArray.Transaction()
						for _, e := range ei {
							t.SaveComponent(e, camera.NewCamera(ecs.GetComponentType(camera.Perspective{})))
						}
						logger.Warn(ecs.FlushMany(t))
					})

					events.Listen(w.EventsBuilder(), func(e sdl.WindowEvent) {
						if e.Event == sdl.WINDOWEVENT_RESIZED {
							events.Emit(w.Events(), camera.NewUpdateProjectionsEvent())
						}
					})
					return nil
				}),
			)
			return nil
		})
	})
}

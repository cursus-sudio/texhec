package projection

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) CameraUp { return CameraUp(mgl32.Vec3{0, 1, 0}) })
	ioc.RegisterSingleton(b, func(c ioc.Dic) CameraForward { return CameraForward(mgl32.Vec3{0, 0, -1}) })

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s cameras.CameraResolverFactory) cameras.CameraResolverFactory {
		getCameraTransform := func(transformArray ecs.ComponentsArray[transform.Transform], entity ecs.EntityID) transform.Transform {
			t, err := transformArray.GetComponent(entity)
			if err != nil {
				return transform.NewTransform()
			}
			return t
		}

		s.Register(ecs.GetComponentType(Ortho{}), func(world ecs.World) func(entity ecs.EntityID) (cameras.Camera, error) {
			transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
			orthoArray := ecs.GetComponentsArray[Ortho](world.Components())
			return func(entity ecs.EntityID) (cameras.Camera, error) {
				getCameraTransformMatrix := func() mgl32.Mat4 {
					cameraTransform := getCameraTransform(transformArray, entity)

					cameraRotation := cameraTransform.Rotation.Inverse()
					cameraPosition := cameraTransform.Rotation.Rotate(cameraTransform.Pos.Mul(-1))
					return cameraRotation.Mat4().Mul4(mgl32.Translate3D(cameraPosition.X(), cameraPosition.Y(), cameraPosition.Z()))
				}
				getProjection := func() Ortho {
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

				cameraGroups, err := ecs.GetComponent[groups.Groups](world.Components(), entity)
				if err != nil {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameras.NewCamera(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					func(mousePos mgl32.Vec2) collider.Ray {
						return ShootRay(
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

		s.Register(ecs.GetComponentType(Perspective{}), func(world ecs.World) func(entity ecs.EntityID) (cameras.Camera, error) {
			transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
			perspectiveArray := ecs.GetComponentsArray[Perspective](world.Components())

			return func(entity ecs.EntityID) (cameras.Camera, error) {
				getCameraTransformMatrix := func() mgl32.Mat4 {
					cameraTransform := getCameraTransform(transformArray, entity)

					up, forward := ioc.Get[CameraUp](c), ioc.Get[CameraForward](c)
					return mgl32.LookAtV(
						cameraTransform.Pos,
						cameraTransform.Pos.Add(cameraTransform.Rotation.Rotate(mgl32.Vec3(forward))),
						mgl32.Vec3(up),
					)
				}
				getProjection := func() Perspective {
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

				cameraGroups, err := ecs.GetComponent[groups.Groups](world.Components(), entity)
				if err != nil {
					cameraGroups = groups.DefaultGroups()
				}
				camera := cameras.NewCamera(
					func() mgl32.Mat4 {
						projMatrix := getProjectionMatrix()
						cameraTransformMatrix := getCameraTransformMatrix()
						return projMatrix.Mul4(cameraTransformMatrix)
					},
					func(mousePos mgl32.Vec2) collider.Ray {
						cameraTransform := getCameraTransform(transformArray, entity)
						return ShootRay(
							getProjectionMatrix(),
							getCameraTransformMatrix(),
							mousePos,
							&cameraTransform.Pos,
						)
					},
					cameraGroups,
				)
				return camera, nil
			}
		})

		return s
	})

}

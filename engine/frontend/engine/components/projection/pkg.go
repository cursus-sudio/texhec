package projection

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/graphics/camera"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s cameras.CameraConstructorsFactory) cameras.CameraConstructorsFactory {
		getCameraTransform := func(transformArray ecs.ComponentsArray[transform.Transform], entity ecs.EntityID) transform.Transform {
			t, err := transformArray.GetComponent(entity)
			if err != nil {
				return transform.NewTransform()
			}
			return t
		}

		s.Register(ecs.GetComponentType(Ortho{}), func(entity ecs.EntityID) (cameras.Camera, error) {
			world := ioc.Get[ecs.World](c)
			transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
			orthoArray := ecs.GetComponentsArray[Ortho](world.Components())

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
			camera := cameras.NewCamera(
				func() mgl32.Mat4 {
					projMatrix := getProjectionMatrix()
					cameraTransformMatrix := getCameraTransformMatrix()
					return projMatrix.Mul4(cameraTransformMatrix)
				},
				func(mousePos mgl32.Vec2) collider.Ray {
					return ShootRay(
						getCameraTransformMatrix(),
						getProjectionMatrix(),
						mousePos,
						nil,
					)
				},
			)
			return camera, nil
		})

		//

		s.Register(ecs.GetComponentType(Perspective{}), func(entity ecs.EntityID) (cameras.Camera, error) {
			world := ioc.Get[ecs.World](c)
			transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
			perspectiveArray := ecs.GetComponentsArray[Perspective](world.Components())

			getCameraTransformMatrix := func() mgl32.Mat4 {
				cameraTransform := getCameraTransform(transformArray, entity)

				return mgl32.LookAtV(
					cameraTransform.Pos,
					cameraTransform.Pos.Add(cameraTransform.Rotation.Rotate(camera.Forward)),
					camera.Up,
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

			camera := cameras.NewCamera(
				func() mgl32.Mat4 {
					projMatrix := getProjectionMatrix()
					cameraTransformMatrix := getCameraTransformMatrix()
					return projMatrix.Mul4(cameraTransformMatrix)
				},
				func(mousePos mgl32.Vec2) collider.Ray {
					cameraTransform := getCameraTransform(transformArray, entity)
					return ShootRay(
						getCameraTransformMatrix(),
						getProjectionMatrix(),
						mousePos,
						&cameraTransform.Pos,
					)
				},
			)
			return camera, nil
		})

		return s
	})

}

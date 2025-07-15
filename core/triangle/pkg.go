package triangle

import (
	"frontend/services/frames"
	"frontend/services/media/window"
	appruntime "shared/services/runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	tools, err := NewTriangleTools()
	if err != nil {
		panic(err.Error())
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) *triangleTools { return tools })

	var t time.Duration
	if false {
		print(t)
	}

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b events.Builder) events.Builder {
		window := ioc.Get[window.Api](c).Window()
		events.Listen(b, func(e frames.FrameEvent) {
			t += e.Delta

			transformSize := [3]float32{100, 100, 0}
			meshSize := [3]float32{100, 100, 0}

			tools.Program.Draw(func() {
				tools.Texture.Draw(func() {
					width, height := window.GetSize()
					{
						gl.Uniform3f(tools.Locations.Resolution, float32(width), float32(height), 1)
					}
					{
						transformSize := [3]float32{
							transformSize[0],
							transformSize[1] * (1 + float32(t.Seconds())),
							transformSize[2],
						}
						scale := [3]float32{
							transformSize[0] / max(1, meshSize[0]),
							transformSize[1] / max(1, meshSize[1]),
							transformSize[2] / max(1, meshSize[2]),
						}
						radians := mgl32.DegToRad(float32(t.Seconds()) * 100)
						rotation := mgl32.QuatIdent().
							Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 0, 1}))
						matrices := []mgl32.Mat4{
							rotation.Mat4(),
							mgl32.Translate3D(
								0-transformSize[0]/2,
								0-transformSize[1]/2,
								0-transformSize[2]/2),
							mgl32.Scale3D(scale[0], scale[1], scale[2]),
						}
						var model mgl32.Mat4
						for i, matrix := range matrices {
							if i == 0 {
								model = matrix
								continue
							}
							model = model.Mul4(matrix)
						}
						gl.UniformMatrix4fv(tools.Locations.Model, 1, false, &model[0])
					}
					{
						camera := mgl32.Translate3D(-0, -0, 0)
						gl.UniformMatrix4fv(tools.Locations.Camera, 1, false, &camera[0])
					}
					{
						projection := mgl32.Ortho2D(
							-float32(width)/2,
							float32(width)/2,
							-float32(height)/2,
							float32(height)/2,
						)
						gl.UniformMatrix4fv(tools.Locations.Projection, 1, false, &projection[0])
					}

					tools.VAO.Draw()
				})
			})
		})
		return b
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		tools := ioc.Get[*triangleTools](c)
		b.OnStop(func(r appruntime.Runtime) {
			tools.Program.Release()
			tools.VAO.Release()
		})
		return b
	})
}

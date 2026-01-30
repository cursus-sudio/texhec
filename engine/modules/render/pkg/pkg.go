package renderpkg

import (
	"bytes"
	"engine/modules/render"
	"engine/modules/render/internal/instancing"
	"engine/modules/render/internal/service"
	"engine/modules/render/internal/systems"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/assets"
	"engine/services/ecs"
	gtexture "engine/services/graphics/texture"
	"engine/services/graphics/vao/vbo"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	for _, pkg := range []ioc.Pkg{
		transitionpkg.PackageT[render.ColorComponent](),
		transitionpkg.PackageT[render.TextureFrameComponent](),
	} {
		pkg.Register(b)
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[render.Vertex] {
		return func() vbo.VBOSetter[render.Vertex] {
			vbo := vbo.NewVBO[render.Vertex](func() {
				gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
					int32(unsafe.Sizeof(render.Vertex{})), uintptr(unsafe.Offsetof(render.Vertex{}.Pos)))
				gl.EnableVertexAttribArray(0)

				gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
					int32(unsafe.Sizeof(render.Vertex{})), uintptr(unsafe.Offsetof(render.Vertex{}.TexturePos)))
				gl.EnableVertexAttribArray(1)
			})
			return vbo
		}
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.Service {
		return service.NewService(c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func() error {
			errs := ecs.RegisterSystems(
				systems.NewClearSystem(c),
				systems.NewErrorLogger(c),
				systems.NewRenderSystem(c),
			)
			if len(errs) != 0 {
				return errs[0]
			}
			return nil
		})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.SystemRenderer {
		return ecs.NewSystemRegister(func() error {
			errs := ecs.RegisterSystems(
				instancing.NewSystem(c),
			)
			if len(errs) != 0 {
				return errs[0]
			}
			return nil
		})
	})

	ioc.WrapService(b, func(c ioc.Dic, b assets.AssetsStorageBuilder) {
		imageHandler := func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			imgFile := bytes.NewBuffer(source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}

			img = gtexture.NewImage(img).FlipV().TrimTransparentBackground().Image()
			return render.NewTextureStorageAsset(img)
		}
		b.RegisterExtension("png", imageHandler)
		b.RegisterExtension("jpg", imageHandler)
		b.RegisterExtension("jpeg", imageHandler)

		b.RegisterExtension("gif", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			imgFile := bytes.NewBuffer(source)
			gif, err := gif.DecodeAll(imgFile)
			if err != nil {
				return nil, err
			}

			images := make([]image.Image, 0, len(gif.Image))
			for _, img := range gif.Image {
				finalImg := gtexture.NewImage(img).FlipV().TrimTransparentBackground().Image()
				images = append(images, finalImg)
			}

			return render.NewTextureStorageAsset(images...)
		})
	})
}

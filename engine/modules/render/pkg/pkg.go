package renderpkg

import (
	"bytes"
	"engine/modules/render"
	"engine/modules/render/internal"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/assets"
	"engine/services/ecs"
	gtexture "engine/services/graphics/texture"
	"image"
	"os"
	"strings"

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

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.Service {
		return internal.NewService(c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func() error {
			ecs.RegisterSystems(
				internal.NewClearSystem(c),
				internal.NewErrorLogger(c),
				internal.NewRenderSystem(c),
			)
			return nil
		})
	})

	ioc.WrapService(b, func(c ioc.Dic, b assets.AssetsStorageBuilder) {
		b.RegisterExtension("png", func(id assets.AssetID) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			imgFile := bytes.NewBuffer(source)
			img, _, err := image.Decode(imgFile)
			if err != nil {
				return nil, err
			}

			img = gtexture.FlipImage(img)
			if !strings.Contains(string(id), "tiles") {
				img = gtexture.TrimTransparentBackground(img)
			}
			return render.NewTextureStorageAsset(img)
		})
	})
}

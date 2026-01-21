package tilepkg

import (
	"bytes"
	"core/modules/tile"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tileservice"
	"core/modules/tile/internal/tileui"
	gridpkg "engine/modules/grid/pkg"
	"engine/services/assets"
	"engine/services/codec"
	"engine/services/ecs"
	gtexture "engine/services/graphics/texture"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	pkgs []ioc.Pkg
}

func Package() ioc.Pkg {
	return pkg{
		[]ioc.Pkg{
			gridpkg.Package[tile.Type](tile.NewTileClickEvent),
			tileservice.Package(),
			tilerenderer.Package(),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}

	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// events
			Register(tile.TileClickEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.System {
		systems := []tile.System{
			tileui.NewSystem(c),
		}
		return ecs.NewSystemRegister(func() error {
			for _, system := range systems {
				if err := system.Register(); err != nil {
					return err
				}
			}
			return nil
		})
	})

	ioc.WrapService(b, func(c ioc.Dic, b assets.AssetsStorageBuilder) {
		b.RegisterExtension("biom", func(id assets.AssetID) (any, error) {
			images := [5]image.Image{}
			directory, _ := strings.CutSuffix(string(id), ".biom")
			for i := range 5 {
				file := fmt.Sprintf("%v/%v.png", directory, i)
				source, err := os.ReadFile(file)
				if err != nil {
					return nil, err
				}
				imgFile := bytes.NewBuffer(source)
				img, _, err := image.Decode(imgFile)
				if err != nil {
					return nil, err
				}
				img = gtexture.FlipImage(img)
				images[i] = img
			}

			return tile.NewBiomAsset(images)
		})
	})
}

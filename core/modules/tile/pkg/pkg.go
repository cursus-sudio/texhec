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
			gridpkg.Package[tile.ID](tile.NewTileClickEvent),
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
			images := [6][]image.Image{}
			directory, _ := strings.CutSuffix(string(id), ".biom")

			for i := range 6 {
				tileDir := fmt.Sprintf("%v/%v", directory, i+1)
				files, err := os.ReadDir(tileDir)
				if err != nil {
					return nil, err
				}
				if len(files) == 0 {
					return nil, fmt.Errorf("there is no tile variant for %v tile", i)
				}

				for _, file := range files {
					filePath := fmt.Sprintf("%v/%v", tileDir, file.Name())
					source, err := os.ReadFile(filePath)
					if err != nil {
						return nil, err
					}
					imgFile := bytes.NewBuffer(source)
					img, _, err := image.Decode(imgFile)
					if err != nil {
						return nil, err
					}
					img = gtexture.NewImage(img).FlipV().Image()
					images[i] = append(images[i], img)
				}
			}

			return tile.NewBiomAsset(images)
		})
	})
}

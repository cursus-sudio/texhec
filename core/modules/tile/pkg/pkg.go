package tilepkg

import (
	"core/modules/tile"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tileservice"
	"core/modules/tile/internal/tileui"
	gridpkg "engine/modules/grid/pkg"
	"engine/services/codec"
	"engine/services/ecs"

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
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}

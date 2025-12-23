package tilepkg

import (
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tiletool"
	"core/modules/tile/internal/tileui"
	"engine/modules/groups"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	tileSize int32,
	gridDepth float32,
	tileGroups groups.GroupsComponent,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ, maxZ int32,
) ioc.Pkg {
	return pkg{
		[]ioc.Pkg{
			tiletool.Package(
				tileSize,
				gridDepth,
				mainLayer,
				layers,
				minX, maxX, minY, maxY, minZ,
			),
			tilecollider.Package(
				tileSize,
				gridDepth,
				tileGroups,
				mainLayer,
				layers,
				minX, maxX, minY, maxY, minZ,
			),
			tilerenderer.Package(
				tileSize,
				gridDepth,
				maxZ-minZ,
				tileGroups,
			),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// types
			Register(tile.Layer(0)).
			Register(tile.ColliderPos{}).
			// events
			Register(tile.TileClickEvent{}).
			// components
			Register(tile.PosComponent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.System {
		systems := []tile.System{
			tileui.NewSystem(
				ioc.Get[logger.Logger](c),
				ioc.Get[tile.ToolFactory](c),
			),
		}
		return ecs.NewSystemRegister(func(world tile.World) error {
			for _, system := range systems {
				if err := system.Register(world); err != nil {
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

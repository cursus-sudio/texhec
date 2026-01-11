package tilepkg

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tileservice"
	"core/modules/tile/internal/tileui"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/services/codec"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	pkgs []ioc.Pkg

	tileSize   int32
	gridDepth  float32
	tileGroups groups.GroupsComponent
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
			tileservice.Package(
				tileSize,
				gridDepth,
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
		tileSize,
		gridDepth,
		tileGroups,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// types
			Register(tile.Layer(0)).
			// events
			Register(tile.TileClickEvent{}).
			// components
			Register(tile.PosComponent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.System {
		systems := []tile.System{
			tileui.NewSystem(c),
			tilecollider.TileColliderSystem(c,
				pkg.tileSize,
				pkg.gridDepth,
				pkg.tileGroups,
				collider.NewCollider(ioc.Get[gameassets.GameAssets](c).SquareCollider),
			),
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

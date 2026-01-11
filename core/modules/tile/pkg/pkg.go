package tilepkg

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tiletool"
	"core/modules/tile/internal/tileui"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/uuid"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

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
			tiletool.Package(
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
		tileToolFactory := ioc.Get[tile.ToolFactory](c)
		logger := ioc.Get[logger.Logger](c)
		systems := []tile.System{
			tileui.NewSystem(
				logger,
				ioc.Get[tile.ToolFactory](c),
			),
			tilecollider.TileColliderSystem(
				tileToolFactory,
				logger,
				pkg.tileSize,
				pkg.gridDepth,
				pkg.tileGroups,
				collider.NewCollider(ioc.Get[gameassets.GameAssets](c).SquareCollider),
				ioc.Get[uuid.Factory](c),
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

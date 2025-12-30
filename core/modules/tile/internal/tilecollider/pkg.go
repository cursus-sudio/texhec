package tilecollider

import (
	gameassets "core/assets"
	"core/modules/tile"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tileSize                     int32
	gridDepth                    float32
	tileGroups                   groups.GroupsComponent
	mainLayer                    tile.Layer
	layers                       []tile.Layer
	minX, maxX, minY, maxY, minZ int32
}

func Package(
	tileSize int32,
	gridDepth float32,
	tileGroups groups.GroupsComponent,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ int32,
) ioc.Pkg {
	return pkg{
		tileSize,
		gridDepth,
		tileGroups,
		mainLayer,
		layers,
		minX, maxX, minY, maxY, minZ,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.System) tile.System {
		tileToolFactory := ioc.Get[tile.ToolFactory](c)
		logger := ioc.Get[logger.Logger](c)
		return ecs.NewSystemRegister(func(w tile.World) error {
			if err := s.Register(w); err != nil {
				return err
			}
			errs := ecs.RegisterSystems(w,
				TileColliderSystem(
					tileToolFactory,
					logger,
					pkg.tileSize,
					pkg.gridDepth,
					pkg.tileGroups,
					collider.NewCollider(ioc.Get[gameassets.GameAssets](c).SquareCollider),
					ioc.Get[uuid.Factory](c),
				),
			)
			if len(errs) != 0 {
				return errs[0]
			}
			return nil
		})
	})
}

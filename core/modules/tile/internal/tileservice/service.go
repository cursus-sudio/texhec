package tileservice

import (
	"core/modules/tile"
	"engine/modules/grid"
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	grid.Service[tile.ID] `inject:"1"`
}

func NewService(c ioc.Dic) tile.Service {
	s := ioc.GetServices[*service](c)
	return s
}

func (t *service) Grid() ecs.ComponentsArray[grid.SquareGridComponent[tile.ID]] {
	return t.Component()
}

func (t *service) GetPos(coords grid.Coords) transform.PosComponent {
	size := t.GetTileSize().Size
	return transform.NewPos(
		size.X()*(float32(coords.X)+.5),
		size.Y()*(float32(coords.Y)+.5),
		size.Z(),
	)
}
func (t *service) GetTileSize() transform.SizeComponent {
	return transform.NewSize(100, 100, 1)
}

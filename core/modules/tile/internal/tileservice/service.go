package tileservice

import (
	"core/modules/tile"
	"engine/modules/grid"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	grid.Service[tile.Type] `inject:"1"`
}

func NewService(c ioc.Dic) tile.Service {
	s := ioc.GetServices[*service](c)
	return s
}

func (t *service) Grid() ecs.ComponentsArray[grid.SquareGridComponent[tile.Type]] {
	return t.Component()
}

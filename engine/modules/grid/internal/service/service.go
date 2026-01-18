package service

import (
	"engine/modules/grid"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type service[Tile grid.TileConstraint] struct {
	World     ecs.World     `inject:"1"`
	Logger    logger.Logger `inject:"1"`
	component ecs.ComponentsArray[grid.SquareGridComponent[Tile]]
}

func NewService[Tile grid.TileConstraint](c ioc.Dic) grid.Service[Tile] {
	s := ioc.GetServices[*service[Tile]](c)
	s.component = ecs.GetComponentsArray[grid.SquareGridComponent[Tile]](s.World)

	return s
}

func (s *service[Tile]) Component() ecs.ComponentsArray[grid.SquareGridComponent[Tile]] {
	return s.component
}

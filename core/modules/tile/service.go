package tile

import (
	"engine/modules/grid"
	"engine/services/ecs"
)

type Service interface {
	Grid() ecs.ComponentsArray[grid.SquareGridComponent[Type]]
}

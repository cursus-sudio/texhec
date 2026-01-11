package tile

import (
	"engine/modules/relation"
	"engine/services/ecs"
)

type Service interface {
	PosKey() relation.Service[PosComponent]

	Pos() ecs.ComponentsArray[PosComponent]
}

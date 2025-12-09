package relation

import (
	"engine/services/ecs"
)

type EntityToKeyTool[Key any] interface {
	Get(Key) (ecs.EntityID, bool)
}

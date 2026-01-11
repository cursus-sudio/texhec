package relation

import (
	"engine/services/ecs"
)

type Service[Key any] interface {
	Get(Key) (ecs.EntityID, bool)
}

package relation

import (
	"engine/services/ecs"
)

type ToolFactory[Key any] ecs.ToolFactory[World, EntityToKeyTool[Key]]
type EntityToKeyTool[Key any] interface {
	Get(Key) (ecs.EntityID, bool)
}

type World interface {
	ecs.World
}

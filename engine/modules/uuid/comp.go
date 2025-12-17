package uuid

import (
	"engine/services/ecs"
)

type Component struct {
	ID UUID
}

func New(id UUID) Component {
	return Component{id}
}

// add unique id to domain components
type UUIDTool interface {
	Factory
	Entity(UUID) (ecs.EntityID, bool)
}

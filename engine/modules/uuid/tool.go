package uuid

import (
	"engine/services/ecs"

	"github.com/google/uuid"
)

// engine interface

type Component struct {
	ID UUID
}

func New(id UUID) Component {
	return Component{id}
}

// add unique id to domain components
type UUIDTool interface {
	UUID() Interface
}

type Interface interface {
	Factory
	Entity(UUID) (ecs.EntityID, bool)
}

// raw interface

type UUID uuid.UUID

func (id UUID) String() string  { return uuid.UUID(id).String() }
func (uuid UUID) Bytes() []byte { return uuid[:] }

type Factory interface {
	NewUUID() UUID
}

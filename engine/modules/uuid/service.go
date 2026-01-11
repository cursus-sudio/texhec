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

type Service interface {
	Factory
	Component() ecs.ComponentsArray[Component]
	Entity(UUID) (ecs.EntityID, bool)
}

// raw interface

type UUID uuid.UUID

func (id *UUID) String() string  { return uuid.UUID(*id).String() }
func (uuid *UUID) Bytes() []byte { return uuid[:] }

type Factory interface {
	NewUUID() UUID
}

package uuid

import (
	"github.com/google/uuid"
)

// interface

type UUID uuid.UUID

func (id UUID) String() string  { return uuid.UUID(id).String() }
func (uuid UUID) Bytes() []byte { return uuid[:] }

type Factory interface {
	NewUUID() UUID
}

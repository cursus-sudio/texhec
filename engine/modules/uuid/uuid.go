package uuid

import (
	"github.com/google/uuid"
)

// interface

type UUID uuid.UUID

func (uuid UUID) String() string { return uuid.String() }
func (uuid UUID) Bytes() []byte  { return uuid[:] }

type Factory interface {
	NewUUID() UUID
}

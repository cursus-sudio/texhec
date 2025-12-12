package servertypes

import (
	"engine/modules/netsync/internal/state"
	"engine/modules/uuid"
)

// server messages

type SendStateDTO struct {
	State state.State
	Error error
}

type SendChangeDTO struct {
	EventID uuid.UUID
	Changes state.State
	Error   error
}

type TransparentEventDTO struct {
	Event any
	Error error
}

package servertypes

import (
	"engine/modules/sync/internal/state"
	"engine/modules/uuid"
)

// server messages

type SendStateDTO struct {
	State state.State
}

type SendChangeDTO struct {
	EventID uuid.UUID
	Changes state.State
}

type TransparentEventDTO struct {
	Event any
}

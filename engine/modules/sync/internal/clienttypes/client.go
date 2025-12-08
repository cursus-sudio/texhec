package clienttypes

import (
	"engine/modules/uuid"
)

// types

type PredictedEvent struct {
	ID    uuid.UUID
	Event any
}

// client messages

type FetchStateDTO struct{}
type EmitEventDTO PredictedEvent

type TransparentEventDTO struct {
	Event any
}

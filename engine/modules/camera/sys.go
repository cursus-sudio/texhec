package camera

import "engine/services/ecs"

// TODO change UpdateProjectionsEvent to some media.ChangedResolution

// updates dynamic projections
type ChangedResolutionEvent struct{}

func NewUpdateProjectionsEvent() ChangedResolutionEvent {
	return ChangedResolutionEvent{}
}

type System ecs.SystemRegister

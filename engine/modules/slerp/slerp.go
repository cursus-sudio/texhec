package slerp

import (
	"engine/services/ecs"
	"time"
)

type System ecs.SystemRegister[World]

type World interface {
	ecs.World
}

//

type Progress float64

type SlerpComponent[Component any] struct {
	From, To Component
	Progress,
	Duration time.Duration
}

func NewSlerp[Component any](from, to Component, duration time.Duration) SlerpComponent[Component] {
	return SlerpComponent[Component]{
		From:     from,
		To:       to,
		Progress: 0,
		Duration: duration,
	}
}

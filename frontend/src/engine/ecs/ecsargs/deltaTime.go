package ecsargs

import "time"

type DeltaTime interface {
	Duration() time.Duration
}

type deltaTime struct{ duration time.Duration }

func NewDeltaTime(duration time.Duration) DeltaTime {
	return &deltaTime{
		duration: duration}
}

func (deltaTime *deltaTime) Duration() time.Duration {
	return deltaTime.duration
}

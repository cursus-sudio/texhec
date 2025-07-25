package frames

import (
	"shared/services/clock"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

const (
	// before scene listeners are executed
	HandleInputs ioc.Order = iota
	FlushDraw
)

type FrameEvent struct {
	Delta time.Duration
}

func NewFrameEvent(delta time.Duration) FrameEvent {
	return FrameEvent{Delta: delta}
}

type Builder interface {
	FPS(int) Builder
	Build(events events.Events, clock clock.Clock) Frames
}

type builder struct {
	fps    int
	clock  clock.Clock
	events events.Events
}

func NewBuilder(fps int) Builder {
	return &builder{
		fps: fps,
	}
}

func (b *builder) FPS(fps int) Builder {
	b.fps = fps
	return b
}

func (b *builder) Build(events events.Events, clock clock.Clock) Frames {
	return &frames{
		Running: false,
		FPS:     b.fps,
		Events:  events,
		Clock:   clock,
	}
}

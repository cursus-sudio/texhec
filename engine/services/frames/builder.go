package frames

import (
	"engine/services/clock"
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

// tick has fixed delta separate from frame rate (usually 30 tick per second)
// tick is triggered before frame as many times as many ticks passed between frames
type TickEvent struct {
	Delta time.Duration
}

func NewFrameEvent(delta time.Duration) FrameEvent {
	return FrameEvent{Delta: delta}
}

type Builder interface {
	// frame per second
	FPS(int) Builder
	// tick per second
	TPS(int) Builder
	Build(events events.Events, clock clock.Clock) Frames
}

type builder struct {
	tps,
	fps int
}

func NewBuilder(tps, fps int) Builder {
	return &builder{
		tps: tps,
		fps: fps,
	}
}

func (b *builder) FPS(fps int) Builder {
	b.fps = fps
	return b
}

func (b *builder) TPS(tps int) Builder {
	b.tps = tps
	return b
}

func (b *builder) Build(events events.Events, clock clock.Clock) Frames {
	return &frames{
		Running: false,
		TPS:     b.tps,
		FPS:     b.fps,
		Events:  events,
		Clock:   clock,
	}
}

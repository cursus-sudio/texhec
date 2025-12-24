package frames

import (
	"engine/services/clock"
	"errors"
	"time"

	"github.com/ogiusek/events"
)

var (
	ErrAlreadyRunning error = errors.New("already running")
)

type Frames interface {
	Run() error
	Stop()
}

type frames struct {
	Running bool
	TPS,
	FPS int
	TickProgress time.Duration
	Events       events.Events
	Clock        clock.Clock
}

func (frames *frames) StartLoop() {
	frameDuration := time.Second / time.Duration(frames.FPS)
	tickDuration := time.Second / time.Duration(frames.TPS)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	var lastFrameTime time.Time
	lastFrameTime = frames.Clock.Now()

	for frames.Running {
		<-ticker.C
		currentTime := time.Now()

		delta := currentTime.Sub(lastFrameTime)
		event := NewFrameEvent(delta)
		frames.TickProgress += delta
		for frames.TickProgress > tickDuration {
			frames.TickProgress -= tickDuration
			events.Emit(frames.Events, TickEvent{tickDuration})
		}
		events.Emit(frames.Events, event)

		lastFrameTime = currentTime
	}
}

func (frames *frames) Run() error {
	if frames.Running {
		return ErrAlreadyRunning
	}

	frames.Running = true
	frames.StartLoop()

	return nil
}

func (frames *frames) Stop() {
	frames.Running = false
}

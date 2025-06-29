package frames

import (
	"errors"
	"shared/services/clock"
	"time"
)

var (
	ErrAlreadyRunning error = errors.New("already running")
)

type Frames interface {
	Run() error
}

type frames struct {
	AlreadyRunning bool
	FPS            int
	OnFrame        func(OnFrame)
	Clock          clock.Clock
}

func (frames *frames) StartLoop() {
	var previousFrame time.Time
	// this should be in game loop when framerate could be changed
	timePerFrame := time.Second / time.Duration(frames.FPS)
	previousFrame = frames.Clock.Now()

	for { // runnning game loop
		now := frames.Clock.Now()

		deltaTime := now.Sub(previousFrame)
		onFrame := NewOnFrame(deltaTime)
		frames.OnFrame(onFrame)

		previousFrame = now
		time.Sleep(previousFrame.Add(timePerFrame).Sub(now))
	}
}

func (frames *frames) Run() error {
	if frames.AlreadyRunning {
		return ErrAlreadyRunning
	}

	go frames.StartLoop()

	return nil
}

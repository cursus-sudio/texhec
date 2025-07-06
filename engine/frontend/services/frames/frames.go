package frames

import (
	"errors"
	"shared/services/clock"
	"sync"
	"time"
)

var (
	ErrAlreadyRunning error = errors.New("already running")
)

type Frames interface {
	Run() error
	Stop()
}

type frames struct {
	AlreadyRunning bool
	Running        bool
	FPS            int
	RunMutex       sync.RWMutex
	OnFrame        func(OnFrame)
	Clock          clock.Clock
}

func (frames *frames) StartLoop() {
	var previousFrame time.Time
	// this should be in game loop when framerate could be changed
	timePerFrame := time.Second / time.Duration(frames.FPS)
	previousFrame = frames.Clock.Now()

	frames.RunMutex.Lock()
	for frames.Running { // runnning game loop
		now := frames.Clock.Now()

		deltaTime := now.Sub(previousFrame)
		onFrame := NewOnFrame(deltaTime)
		frames.OnFrame(onFrame)

		previousFrame = now
		time.Sleep(previousFrame.Add(timePerFrame).Sub(now))
	}
	frames.RunMutex.Unlock()
}

func (frames *frames) Run() error {
	if frames.AlreadyRunning {
		return ErrAlreadyRunning
	}

	frames.StartLoop()

	return nil
}

func (frames *frames) Stop() {
	frames.Running = false
	frames.RunMutex.RLock()
	frames.RunMutex.RUnlock()
}

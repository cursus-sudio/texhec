package frames

import (
	"shared/services/clock"
	"time"

	"github.com/ogiusek/ioc/v2"
)

const (
	HandleInputs ioc.Order = iota
	Update
	BeforeDraw
	Clear
	Draw
	AfterDraw
)

type OnFrame struct {
	Delta time.Duration
}

func NewOnFrame(delta time.Duration) OnFrame {
	return OnFrame{
		Delta: delta,
	}
}

type Builder interface {
	FPS(int) Builder
	OnFrame(func(OnFrame)) Builder
	Build() Frames
}

type builder struct {
	fps     int
	onFrame []func(OnFrame)
	clock   clock.Clock
}

func NewBuilder(clock clock.Clock) Builder {
	return &builder{
		fps:     60,
		onFrame: []func(OnFrame){},
		clock:   clock,
	}
}

func (b *builder) FPS(fps int) Builder {
	b.fps = fps
	return b
}

func (b *builder) OnFrame(onFrame func(OnFrame)) Builder {
	b.onFrame = append(b.onFrame, onFrame)
	return b
}

func (b *builder) Build() Frames {
	return &frames{
		AlreadyRunning: false,
		FPS:            b.fps,
		Running:        true,
		OnFrame: func(of OnFrame) {
			for _, onFrame := range b.onFrame {
				onFrame(of)
			}
		},
		Clock: b.clock,
	}
}

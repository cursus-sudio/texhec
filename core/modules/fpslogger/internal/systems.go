package internal

import (
	"engine/services/console"
	"engine/services/ecs"
	"engine/services/frames"
	"fmt"
	"time"

	"github.com/ogiusek/events"
)

type logsSystem struct {
	World   ecs.World
	Console console.Console

	frames []time.Time
}

func NewFpsLoggerSystem(
	console console.Console,
) ecs.SystemRegister[ecs.World] {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &logsSystem{
			World:   w,
			Console: console,

			frames: make([]time.Time, 60),
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

var format = "02-01-2006 15:04:05"

func (system *logsSystem) Listen(args frames.FrameEvent) {
	now := time.Now()
	latestAcceptableFrame := now.Add(-time.Second)
	startIndex := 0
	for i, frame := range system.frames {
		if latestAcceptableFrame.Before(frame) {
			startIndex = i
			break
		}
	}
	system.frames = append(system.frames[startIndex:], now)

	//

	text := "----------------------------------------------------------------\n"
	text += fmt.Sprintf("now %s\n", time.Now().Format(format))
	text += fmt.Sprintf("fps %d\n", len(system.frames))

	system.Console.Print(text)
	system.Console.Flush()
}

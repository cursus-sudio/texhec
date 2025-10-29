package fpslogger

import (
	"fmt"
	"frontend/services/console"
	"frontend/services/frames"
	"frontend/services/scenes"
	"shared/services/ecs"
	"time"

	"github.com/ogiusek/events"
)

type logsSystem struct {
	SceneManager scenes.SceneManager
	World        ecs.World
	Console      console.Console

	frames []time.Time
}

func NewFpsLoggerSystem(
	sceneMagener scenes.SceneManager,
	console console.Console,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &logsSystem{
			SceneManager: sceneMagener,
			World:        w,
			Console:      console,

			frames: make([]time.Time, 60),
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

var format = "02-01-2006 15:04:05"

func (system *logsSystem) Listen(args frames.FrameEvent) error {
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
	return nil
}

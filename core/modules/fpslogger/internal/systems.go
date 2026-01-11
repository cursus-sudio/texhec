package internal

import (
	"engine/services/console"
	"engine/services/ecs"
	"engine/services/frames"
	"fmt"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type logsSystem struct {
	World         ecs.World       `inject:"1"`
	Console       console.Console `inject:"1"`
	EventsBuilder events.Builder  `inject:"1"`

	frames []time.Time
}

func NewFpsLoggerSystem(c ioc.Dic) ecs.SystemRegister {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*logsSystem](c)
		s.frames = make([]time.Time, 60)
		events.Listen(s.EventsBuilder, s.Listen)
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

package logs

import (
	"fmt"
	"frontend/services/console"
	"frontend/services/frames"
	"frontend/services/scenes"
	"shared/services/ecs"
	"sync"
	"time"

	"github.com/ogiusek/events"
)

type logsSystem struct {
	SceneManager scenes.SceneManager
	World        ecs.World
	Console      console.Console

	Mutex sync.Mutex
	Fps   int
}

func NewLogsSystem(
	sceneMagener scenes.SceneManager,
	world ecs.World,
	console console.Console,
) ecs.SystemRegister {
	return &logsSystem{
		SceneManager: sceneMagener,
		World:        world,
		Console:      console,
	}
}

func (s *logsSystem) Register(b events.Builder) {
	events.ListenE(b, s.Listen)
}

var format = "02-01-2006 15:04:05"

func (system *logsSystem) Listen(args frames.FrameEvent) error {
	go func() {
		system.Mutex.Lock()
		system.Fps++
		system.Mutex.Unlock()
		time.Sleep(time.Second)
		system.Mutex.Lock()
		system.Fps--
		system.Mutex.Unlock()
	}()
	text := "----------------------------------------------------------------\n"
	text += fmt.Sprintf("now %s\n", time.Now().Format(format))
	text += fmt.Sprintf("fps %d\n", system.Fps)

	system.Console.Print(text)
	system.Console.Flush()
	return nil
}

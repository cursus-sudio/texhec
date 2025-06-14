package gameloop

import (
	"frontend/src/engine/ecs"
	"frontend/src/engine/ecs/ecsargs"
	"time"

	"github.com/ogiusek/ioc"
)

func StartGameLoop(c ioc.Dic) {
	game := ioc.GetServices[Game](c)
	game.LoadDefaults()

	var previousFrame time.Time
	previousFrame = time.Now()

	for {
		now := time.Now()
		deltaTime := ecsargs.NewDeltaTime(now.Sub(previousFrame))

		game.Update(ecs.NewArgs(deltaTime))

		previousFrame = now
		time.Sleep(previousFrame.Add(time.Millisecond * 16).Sub(now))
	}
}

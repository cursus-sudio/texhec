package gameloop

import (
	"backend/services/api"
	"fmt"
	"frontend/example/tacticalmap"
	"frontend/services/draw"
	"frontend/services/ecs"
	"strings"
	"time"

	"github.com/ogiusek/relay/v2"
)

type Game struct {
	Backend    api.Server `inject:"1"`
	StartedAt  time.Time
	Frame      int
	PastTime   time.Duration
	DrawnLines int
}

func (game *Game) LoadDefaults() {
	game.StartedAt = time.Now()
	game.Frame = 0
	game.PastTime = time.Duration(0)
	game.DrawnLines = 0
}

func (game *Game) Update(args ecs.Args) {
	game.PastTime += args.DeltaTime.Duration()
	game.Frame += 1
}

func goToPreviousLine() {
	print("\033[1A")
}

func clearLine() {
	print("\033[2K")
}

func (game *Game) Draw(args draw.DrawApi) {
	for i := 0; i < game.DrawnLines; i++ {
		goToPreviousLine()
		clearLine()
	}

	format := "02-01-2006 15:04:05"
	text := ""
	text += fmt.Sprintf("first frame %s\n", game.StartedAt.Format(format))
	text += fmt.Sprintf("now %s\n", time.Now().Format(format))
	text += fmt.Sprintf("current frame %d\n", game.Frame)
	text += fmt.Sprintf("time in game %f\n", game.PastTime.Seconds())

	tacticalMap, _ := relay.Handle(game.Backend.Relay(), tacticalmap.NewGetReq())
	text += fmt.Sprintf("found shit %v\n", tacticalMap)

	fmt.Print(text)

	lines := strings.Count(text, "\n")
	game.DrawnLines = lines
}

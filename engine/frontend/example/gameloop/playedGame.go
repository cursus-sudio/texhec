package gameloop

import (
	"backend/services/api"
	"frontend/services/inputs"
	"frontend/services/window"
)

type PlayedGame struct {
	backend api.Server
	// idk received data on this end
}

func (game *PlayedGame) Update(inputs inputs.Inputs) {
}

func (game *PlayedGame) Draw(window window.Window) {
}

func NewPlayedGame(
	backend api.Server,
) PlayedGame {
	return PlayedGame{
		backend: backend,
	}
}

package gameloop

import (
	"backend/services/backendapi"
	"frontend/services/inputs"
	"frontend/services/window"
)

type PlayedGame struct {
	backend backendapi.Backend
	// idk received data on this end
}

func (game *PlayedGame) Update(inputs inputs.Inputs) {
}

func (game *PlayedGame) Draw(window window.Window) {
}

func NewPlayedGame(
	backend backendapi.Backend,
) PlayedGame {
	return PlayedGame{
		backend: backend,
	}
}

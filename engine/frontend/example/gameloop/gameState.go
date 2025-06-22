package gameloop

import (
	"backend/services/api"
	"frontend/services/inputs"
	"frontend/services/window"

	"github.com/ogiusek/null"
)

type GameState struct {
	// everything
	// save. is game going or paused

	PlayedGame null.Nullable[PlayedGame]
}

func (gameState *GameState) LoadGame(backend api.Server) {
	gameState.PlayedGame = null.New(NewPlayedGame(backend))
}

func (gameState *GameState) Update(inputs inputs.Inputs) {
	if playedGame, ok := gameState.PlayedGame.Ok(); ok {
		// if time plays
		playedGame.Update(inputs)
	}
}

func (gameState *GameState) Draw(window window.Window) {
	if playedGame, ok := gameState.PlayedGame.Ok(); ok {
		playedGame.Draw(window)
		// if paused draw overlay
	}
}

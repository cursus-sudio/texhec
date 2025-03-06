package gameobjects

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
)

type GameState struct {
	models.ModelBase
	players []*Player
	gameMap *Map
}

func NewGame(c ioc.Dic, base models.ModelBase, players []*Player, gameMap *Map) GameState {
	return GameState{
		ModelBase: base,
		players:   players,
		gameMap:   gameMap,
	}
}

// func GetMap()

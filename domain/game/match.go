package game

import (
	"domain/blueprints"
	"domain/common/models"
	domain "domain/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type Game struct {
	models.ModelBase
	Seed    vo.Seed
	Turn    int
	Players []*domain.Player
	GameMap *domain.MapInGame
	Battle  null.Nullable[*Battle]
}

func NewGameState(c ioc.Dic, prepare Hub) *Game {
	usersLen := len(prepare.Users)
	players := make([]*domain.Player, usersLen)
	bioms := make([]*blueprints.Biom, usersLen)
	brains := make(map[models.ModelId]*domain.Building, usersLen)
	for i, user := range prepare.Users {
		fraction := prepare.UserFractions[user.Id]
		biom := fraction.SpawnBiom
		bioms = append(bioms, biom)

		brain := domain.NewBuilding(c, fraction.BrainBlueprint)
		brains[biom.Id] = brain

		player := domain.NewPlayer(c,
			prepare.Users[i],
			fraction,
			brain,
			prepare.Options.StartBudget,
		)
		players = append(players, player)
	}

	gameMap := domain.NewMapInGame(c, prepare.Options.MapSize, bioms, brains)

	return &Game{
		ModelBase: models.NewBase(c),
		Seed:      prepare.Options.Seed,
		Players:   players,
		GameMap:   gameMap,
		Battle:    null.Null[*Battle](),
	}

}

// func (game *GameStateGame)

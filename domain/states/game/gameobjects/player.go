package gameobjects

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/kernel"
	"log"
	"os/user"

	"github.com/ogiusek/ioc"
)

// SERVICE
type PlayerErrors interface {
	DoNotFoundResearch(player Player, research Research) error
	PlayerAlreadyMakesTurn(player Player) error
	PlayerDoesNotMakeTurn(player Player) error
}

//

type Player struct {
	models.ModelBase
	user      *user.User
	makesTurn bool
	fraction  gamedefinitions.Fraction
	unlocked  map[models.ModelId]kernel.ResearchDefinition
	studies   map[models.ModelId]Research
}

func NewPlayer(c ioc.Dic, base models.ModelBase, user *user.User, fraction gamedefinitions.Fraction) Player {
	research := fraction.InitialResearch()
	unlocked := map[models.ModelId]kernel.ResearchDefinition{research.Id(): research}
	studies := map[models.ModelId]Research{}
	for _, definition := range research.NextResearches() {
		// THIS IS NOT CODE
		log.Panic(definition)
		// study := NewResearch(c, models.NewBase(c), definition, vo.EmptyBudget())
		// studies[study.Id()] = study
	}
	return Player{
		ModelBase: base,
		user:      user,
		makesTurn: false,
		fraction:  fraction,
		unlocked:  unlocked,
		studies:   studies,
	}
}

func (player *Player) User() user.User {
	return *player.user
}

func (player *Player) MakesTurn() bool {
	return player.makesTurn
}

func (player *Player) Fraction() gamedefinitions.Fraction {
	return player.fraction
}

func (player *Player) Researched() map[models.ModelId]kernel.ResearchDefinition {
	researched := make(map[models.ModelId]kernel.ResearchDefinition, len(player.unlocked))
	for _, def := range player.unlocked {
		researched[def.Id()] = def
	}
	return researched
}

func (player *Player) Studies() map[models.ModelId]Research {
	studies := make(map[models.ModelId]Research, len(player.studies))
	for _, study := range player.studies {
		studies[study.Id()] = study
	}
	return studies
}

func (player *Player) StartTurn(c ioc.Dic) error {
	if player.makesTurn {
		errors := ioc.Get[PlayerErrors](c)
		return errors.PlayerAlreadyMakesTurn(*player)
	}
	player.makesTurn = true
	return nil
}

func (player *Player) EndTurn(c ioc.Dic) error {
	if !player.makesTurn {
		errors := ioc.Get[PlayerErrors](c)
		return errors.PlayerDoesNotMakeTurn(*player)
	}
	player.makesTurn = false
	return nil
}

func (player *Player) UpdateResearch(c ioc.Dic, research Research) error {
	unlock, ok := player.studies[research.Id()]
	if !ok {
		errors := ioc.Get[PlayerErrors](c)
		return errors.DoNotFoundResearch(*player, research)
	}

	if research.IsComplete() {
		player.unlocked[unlock.definition.Id()] = unlock.definition
		delete(player.studies, unlock.Id())
	} else {
		// player.studies[unlock.Id()].progress = research.progress
	}

	return nil
}

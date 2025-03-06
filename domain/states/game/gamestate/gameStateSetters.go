package gamestate

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/game/gamevo"
	"domain/states/kernel/vo"
)

type FogOfWarSetter interface {
	models.ModelBaseSetter
	SeeTile(...gamevo.Pos) error
	CoverTile(...gamevo.Pos) error
}

type ResearchSetter interface {
	models.ModelBaseSetter
	// automatically changes player
	IncreaseProgress(vo.Budget) error
}

type PlayerSetter interface {
	// StartTurn is run automatically when previous player ends its turn
	EndTurn() error
}

type TroopBattalionSetter interface {
	DecreaseBudget(vo.TroopBudget) error

	IncreaseBrigadeSize(uint) error
	DecreaseBrigadeSize(uint) error

	UseSkill(gamedefinitions.TroopSkill) error
}

type BuildingSetter interface {
	Heal(vo.Hp) error
	Deal(vo.Hp) error
}

type ResourceSetter interface{}

type TileSetter interface {
	Build(owner Player, building gamedefinitions.Building) error
	SpawnBattalion(owner Player, troop gamedefinitions.Troop, amount uint) error

	MoveBattalionsTo([]TroopBattalion, ...gamevo.Direction)

	SettleBattle( /* TODO create type for settled battle */ )
}

type MapSetter interface{}

type GameStateSetter interface {
	models.ModelBaseSetter
}

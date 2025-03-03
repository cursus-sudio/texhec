package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

// SERVICE
type TileInGameErrors struct {
	TileAlreadyHasBuilding,
	CannotBuildOnTile,
	CannotBuildBelowTroop,
	ThisBuildingCannotBeBuildHere error
}

type TileInGame struct {
	models.ModelBase
	Position  vo.PositionInGame
	Blueprint *blueprints.TileInGame
	Resource  null.Nullable[*blueprints.Resource]
	Building  null.Nullable[*Building]
	Troop     null.Nullable[*TroopArmyInGame]
}

func NewTileInGame(c ioc.Dic, position vo.PositionInGame, blueprint *blueprints.TileInGame) *TileInGame {
	return &TileInGame{
		ModelBase: models.NewBase(c),
		Position:  position,
		Blueprint: blueprint,
		Resource:  null.Null[*blueprints.Resource](),
		Building:  null.Null[*Building](),
		Troop:     null.Null[*TroopArmyInGame](),
	}
}

func (tile *TileInGame) CanBuild(c ioc.Dic, blueprint *blueprints.Building) error {
	errors := ioc.Get[TileInGameErrors](c)
	if tile.Building.Ok {
		return errors.TileAlreadyHasBuilding
	}
	if tile.Blueprint.Occupied {
		return errors.CannotBuildOnTile
	}
	if tile.Troop.Ok {
		return errors.CannotBuildBelowTroop
	}
	if !blueprint.CanBuildOn(tile.Resource) {
		return errors.ThisBuildingCannotBeBuildHere
	}
	return nil
}

func (tile *TileInGame) Build(c ioc.Dic, blueprint *blueprints.Building) error {
	if err := tile.CanBuild(c, blueprint); err != nil {
		return err
	}

	building := NewBuilding(c, blueprint)
	tile.Building = null.New(building)
	if tile.Blueprint.OnBuildChangesIn.Ok {
		tile.Blueprint = tile.Blueprint.OnBuildChangesIn.Val
	}

	return nil
}

func (tile *TileInGame) CanMoveHere(c ioc.Dic, blueprint *TroopArmyInGame) error {
	// TODO
	return nil
}

func (tile *TileInGame) MoveHere(c ioc.Dic, blueprint *TroopArmyInGame) error {
	// TODO
	return nil
}

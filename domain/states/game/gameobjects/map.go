package gameobjects

import (
	"domain/common/models"
	"domain/states/game/gamevo"

	"github.com/ogiusek/ioc"
)

type Map struct {
	models.ModelBase
	size  gamevo.Size
	tiles map[gamevo.Pos]*Tile
	// TODO
	// troopsPositions    map[models.ModelId]gamevo.Pos
	// buildingsPositions map[models.ModelId]gamevo.Pos
}

func NewMap(c ioc.Dic, base models.ModelBase, size gamevo.Size, tiles map[gamevo.Pos]*Tile) Map {
	// troopsPositions := make(map[models.ModelId]gamevo.Pos)
	// buildingsPositions := make(map[models.ModelId]gamevo.Pos)
	// for pos, tile := range tiles {
	// 	if tile
	// }

	return Map{
		ModelBase: base,
		size:      size,
		tiles:     tiles,
	}
}

// func GetTile()

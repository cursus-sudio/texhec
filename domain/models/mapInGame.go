package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type MapInGame struct {
	models.ModelBase
	Size  vo.SizeInGame
	Tiles map[vo.PositionInGame]*TileInGame
}

// SERVICE
type MapInGameGenerator interface {
	// TODO
	// idea for this

	// biom generation:
	// pick random points for each player (in this point place brain buildings)
	// create map with distance for each tile +- random offset
	// pick biggest value from each map and mark it as this biom

	// resource generation:
	// determine amount of points by biom size
	// pick points for each biom and grow them to deposit site which is also randomly picked
	New(size vo.SizeInGame, bioms []*blueprints.Biom, brainsByBiomId map[models.ModelId]*Building) *MapInGame
}

func NewMapInGame(c ioc.Dic, size vo.SizeInGame, bioms []*blueprints.Biom, brainsByBiomId map[models.ModelId]*Building) *MapInGame {
	return ioc.Get[MapInGameGenerator](c).New(size, bioms, brainsByBiomId)
}

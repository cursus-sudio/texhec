package blueprints

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type TileInGame struct {
	models.ModelBase
	models.ModelDescription
	Resources        []*TileResource
	OnBuildChangesIn null.Nullable[*TileInGame]
	Occupied         bool
	DebuffedMoveDifficulty,
	BuffedMoveDifficulty uint
}

func NewTileInGame(c ioc.Dic, desc models.ModelDescription, resources []*TileResource, onBuildChangesIn null.Nullable[*TileInGame], occupied bool, debuffedMoveDifficulty, buffedMoveDifficulty uint) TileInGame {
	return TileInGame{
		ModelBase:              models.NewBase(c),
		ModelDescription:       desc,
		Resources:              resources,
		OnBuildChangesIn:       onBuildChangesIn,
		Occupied:               occupied,
		DebuffedMoveDifficulty: debuffedMoveDifficulty,
		BuffedMoveDifficulty:   buffedMoveDifficulty,
	}
}

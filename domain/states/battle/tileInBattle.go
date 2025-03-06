package battle

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type TileType struct {
	models.ModelBase
	models.ModelDescription
	Obstacle null.Nullable[*Obstacle]
	DebuffedMoveDifficulty,
	BuffedMoveDifficulty uint
}

func NewTileInBattle(c ioc.Dic, base models.ModelBase, desc models.ModelDescription, obstacle null.Nullable[*Obstacle], debuffedMoveDifficulty, buffedMoveDifficulty uint) TileType {
	return TileType{
		ModelBase:              base,
		ModelDescription:       desc,
		Obstacle:               obstacle,
		DebuffedMoveDifficulty: debuffedMoveDifficulty,
		BuffedMoveDifficulty:   buffedMoveDifficulty,
	}
}

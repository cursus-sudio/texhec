package blueprints

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
)

type TileInBattle struct {
	models.ModelBase
	models.ModelDescription
	Full bool
	DebuffedMoveDifficulty,
	BuffedMoveDifficulty uint
}

func NewTileInBattle(c ioc.Dic, desc models.ModelDescription, full bool, debuffedMoveDifficulty, buffedMoveDifficulty uint) TileInBattle {
	return TileInBattle{
		ModelBase:              models.NewBase(c),
		ModelDescription:       desc,
		Full:                   full,
		DebuffedMoveDifficulty: debuffedMoveDifficulty,
		BuffedMoveDifficulty:   buffedMoveDifficulty,
	}
}

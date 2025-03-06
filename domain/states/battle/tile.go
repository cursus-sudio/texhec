package battle

import (
	"domain/common/models"
	battlevo "domain/states/battle/vo"

	"github.com/ogiusek/null"
)

type Tile struct {
	models.ModelBase
	Position  battlevo.Pos
	Blueprint *TileType
	Troop     null.Nullable[*Troop]
	Obstacle  null.Nullable[*Obstacle]
}

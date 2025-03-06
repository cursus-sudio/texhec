package battle

import (
	"domain/common/models"
	battlevo "domain/states/battle/vo"
)

type Map struct {
	models.ModelBase
	Size  battlevo.Size
	Tiles map[battlevo.Pos]*Tile
}

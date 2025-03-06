package battle

import "domain/common/models"

type Obstacle struct {
	models.ModelBase
	models.ModelDescription
	Tiles []*TileType
}

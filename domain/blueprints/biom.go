package blueprints

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
)

type Biom struct {
	models.ModelBase
	models.ModelDescription
	Tiles           []*TileInBattle
	SizeMultiplayer float32
}

func NewBiom(c ioc.Dic, desc models.ModelDescription, tiles []*TileInBattle, sizeMultiplayer float32) Biom {
	return Biom{
		ModelBase:        models.NewBase(c),
		ModelDescription: desc,
		Tiles:            tiles,
		SizeMultiplayer:  sizeMultiplayer,
	}
}

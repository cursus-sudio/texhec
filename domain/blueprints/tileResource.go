package blueprints

import "domain/common/models"

type TileResource struct {
	TileId, ResourceId models.ModelId
	SpawnChance        float32
}

func NewTileResource(tileId, resourceId models.ModelId, spawnChange float32) TileResource {
	return TileResource{
		TileId:      tileId,
		ResourceId:  resourceId,
		SpawnChance: spawnChange,
	}
}

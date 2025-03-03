package models

import "domain/common/models"

type DiscoveredTile struct {
	// TODO
	// fog of war should have 3 states.
	// 1. never seen
	// 2. seen
	// 3. sees
	PlayerId, TileId models.ModelId
}

func NewDiscoveredTile(playerId, tileId models.ModelId) DiscoveredTile {
	return DiscoveredTile{
		PlayerId: playerId,
		TileId:   tileId,
	}
}

package registry

import "core/modules/tile"

const (
	// this has to be changed before using saving
	// when i save game and then update to new version which adds new tile as not last
	// than other tiles are pushed back and every next tile is marked as previous
	_ tile.ID = iota
	TileWater
	TileSand
	TileGrass
	TileMountain
)

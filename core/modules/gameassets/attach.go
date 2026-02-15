package gameassets

import (
	"core/modules/construct"
	"core/modules/tile"
	"engine/modules/transition"
	"engine/services/assets"
)

type GameAssets struct {
	Hud   HudAssets
	Tiles TileAssets
	Units UnitAssets

	ExampleAudio assets.AssetID `path:"audio.wav"`

	Blank          assets.AssetID `path:"blank texture"`
	SquareMesh     assets.AssetID `path:"square mesh"`
	SquareCollider assets.AssetID `path:"square collider"`
	FontAsset      assets.AssetID `path:"font1.ttf"`
}

type HudAssets struct {
	Btn         assets.AssetID `path:"hud/btn.png"`
	Cursor      assets.AssetID `path:"hud/cursor.png"`
	Settings    assets.AssetID `path:"hud/settings.png"`
	Background1 assets.AssetID `path:"hud/bg1.gif"`
	Background2 assets.AssetID `path:"hud/bg2.gif"`
}

type TileAssets struct {
	Grass    assets.AssetID `path:"tiles/grass.biom"`
	Sand     assets.AssetID `path:"tiles/sand.biom"`
	Mountain assets.AssetID `path:"tiles/mountain.biom"`
	Water    assets.AssetID `path:"tiles/water.biom"`
}

type UnitAssets struct {
	Unit assets.AssetID `path:"units/tank.png"`
}

//
//
//

const (
	_ transition.EasingID = iota
	LinearEasingFunction
	MyEasingFunction
	EaseOutElastic
)

const (
	_ tile.Type = iota
	TileWater
	TileSand
	TileGrass
	TileMountain
)

const (
	_ construct.ID = iota
	ConstructTank
)

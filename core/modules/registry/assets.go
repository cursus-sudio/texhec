package registry

import "engine/services/assets"

type Assets struct {
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

package registry

import "engine/modules/assets"

// asset ID should be a number.
// asset path and its dispatcher should be pointed by id.
// and this approach should be used by every registry
type Assets struct {
	Hud        HudAssets
	Tiles      TileAssets
	Constructs ConstructAssets
	Units      UnitAssets

	ExampleAudio assets.ID `path:"audio.wav"`

	Blank          assets.ID `path:"blank texture"`
	SquareMesh     assets.ID `path:"square mesh"`
	SquareCollider assets.ID `path:"square collider"`
	FontAsset      assets.ID `path:"font1.ttf"`
}

type HudAssets struct {
	Btn         assets.ID `path:"hud/btn.png-trim"`
	Cursor      assets.ID `path:"hud/cursor.png-trim"`
	Settings    assets.ID `path:"hud/settings.png-trim"`
	Background1 assets.ID `path:"hud/bg1.gif-trim"`
	Background2 assets.ID `path:"hud/bg2.gif-trim"`
}

type TileAssets struct {
	Grass    assets.ID `path:"tiles/grass.biom"`
	Sand     assets.ID `path:"tiles/sand.biom"`
	Mountain assets.ID `path:"tiles/mountain.biom"`
	Water    assets.ID `path:"tiles/water.biom"`
}

type ConstructAssets struct {
	Farm assets.ID `path:"constructs/farm.png"`
}

type UnitAssets struct {
	Tank assets.ID `path:"units/tank.png"`
}

package definitions

import (
	"engine/services/ecs"
)

// asset ID should be a number.
// asset path and its dispatcher should be pointed by id.
// and this approach should be used by every registry
type Definitions struct {
	Hud        Hud
	Tiles      Tiles
	Constructs Constructs
	Units      Units

	Transitions Transitions

	ExampleAudio ecs.EntityID `path:"audio.wav"`

	Blank          ecs.EntityID `path:"blank texture"`
	SquareMesh     ecs.EntityID `path:"square mesh"`
	SquareCollider ecs.EntityID `path:"square collider"`
	FontAsset      ecs.EntityID `path:"font1.ttf"`
}

type Hud struct {
	Btn         ecs.EntityID `path:"hud/btn.png-trim"`
	Cursor      ecs.EntityID `path:"hud/cursor.png-trim"`
	Settings    ecs.EntityID `path:"hud/settings.png-trim"`
	Background1 ecs.EntityID `path:"hud/bg1.gif-trim"`
	Background2 ecs.EntityID `path:"hud/bg2.gif-trim"`
}

type Tiles struct {
	Grass    ecs.EntityID `path:"tiles/grass.biom" tile:""`
	Sand     ecs.EntityID `path:"tiles/sand.biom" tile:""`
	Mountain ecs.EntityID `path:"tiles/mountain.biom" tile:""`
	Water    ecs.EntityID `path:"tiles/water.biom" tile:""`
}

type Constructs struct {
	Farm ecs.EntityID `path:"constructs/farm.png" construct:"farm"`
}

type Units struct {
	Tank ecs.EntityID `path:"units/tank.png"`
}

type Transitions struct {
	Linear         ecs.EntityID `transition:"linear"`
	MyEasing       ecs.EntityID `transition:"my easing"`
	EaseOutElastic ecs.EntityID `transition:"ease out elastic"`
}

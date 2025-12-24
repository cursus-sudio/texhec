package tile

import (
	"core/modules/definition"
	"core/modules/ui"
	"engine"
	"engine/modules/relation"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, TileTool]
type TileTool interface {
	Tile() Interface
}
type World interface {
	engine.World
	definition.DefinitionTool
	ui.UiTool
}
type Interface interface {
	PosKey() relation.EntityToKeyTool[PosComponent]
	ColliderPosKey() relation.EntityToKeyTool[ColliderPos]

	Pos() ecs.ComponentsArray[PosComponent]
}

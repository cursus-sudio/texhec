package definition

import (
	"engine"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, DefinitionTool]
type DefinitionTool interface {
	Definition() Interface
}
type World interface {
	engine.World
}
type Interface interface {
	Component() ecs.ComponentsArray[DefinitionComponent]
	Link() ecs.ComponentsArray[DefinitionLinkComponent]
}

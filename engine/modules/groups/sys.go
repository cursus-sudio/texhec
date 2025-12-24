package groups

import (
	"engine/modules/hierarchy"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, GroupsTool]
type GroupsTool interface {
	Groups() Interface
}
type World interface {
	ecs.World
	hierarchy.HierarchyTool
}
type Interface interface {
	Component() ecs.ComponentsArray[GroupsComponent]
	Inherit() ecs.ComponentsArray[InheritGroupsComponent]
}

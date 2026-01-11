package groups

import (
	"engine/services/ecs"
)

type Service interface {
	Component() ecs.ComponentsArray[GroupsComponent]
	Inherit() ecs.ComponentsArray[InheritGroupsComponent]
}

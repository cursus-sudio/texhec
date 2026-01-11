package internal

import (
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/services/ecs"
	"engine/services/logger"
)

type service struct {
	logger logger.Logger

	world     ecs.World
	hierarchy hierarchy.Service

	inheritArray ecs.ComponentsArray[groups.InheritGroupsComponent]
	groupsArray  ecs.ComponentsArray[groups.GroupsComponent]
}

func NewService(
	w ecs.World,
	hierarchy hierarchy.Service,
	logger logger.Logger,
) groups.Service {
	t := &service{
		logger: logger,

		world:     w,
		hierarchy: hierarchy,

		inheritArray: ecs.GetComponentsArray[groups.InheritGroupsComponent](w),
		groupsArray:  ecs.GetComponentsArray[groups.GroupsComponent](w),
	}
	t.Init()
	return t
}

func (t *service) Component() ecs.ComponentsArray[groups.GroupsComponent] {
	return t.groupsArray
}
func (t *service) Inherit() ecs.ComponentsArray[groups.InheritGroupsComponent] {
	return t.inheritArray
}

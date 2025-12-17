package test

import (
	"engine/modules/groups"
	groupspkg "engine/modules/groups/pkg"
	"engine/modules/hierarchy"
	hierarchypkg "engine/modules/hierarchy/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"testing"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type world struct {
	World         ecs.World
	Hierarchy     hierarchy.Interface
	Groups        ecs.ComponentsArray[groups.GroupsComponent]
	InheritGroups ecs.ComponentsArray[groups.InheritGroupsComponent]
}

type Setup struct {
	world
	T *testing.T
}

func NewSetup(t *testing.T) Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		hierarchypkg.Package(),
		groupspkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	w := ecs.NewWorld()
	ecs.RegisterSystems(w,
		ioc.Get[groups.System](c),
	)
	return Setup{
		world{
			w,
			ioc.Get[ecs.ToolFactory[hierarchy.HierarchyTool]](c).Build(w).Hierarchy(),
			ecs.GetComponentsArray[groups.GroupsComponent](w),
			ecs.GetComponentsArray[groups.InheritGroupsComponent](w),
		},
		t,
	}
}

func (setup Setup) expectGroups(entity ecs.EntityID, expectedGroups groups.GroupsComponent) {
	groups, _ := setup.Groups.Get(entity)
	if groups != expectedGroups {
		setup.T.Errorf("expected pos %v but has %v", expectedGroups, groups)
	}
}

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
	ecs.World
	hierarchy.HierarchyTool
	groups.GroupsTool
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

	w := world{
		World: ecs.NewWorld(),
	}
	w.HierarchyTool = ioc.Get[hierarchy.ToolFactory](c).Build(w)
	w.GroupsTool = ioc.Get[groups.ToolFactory](c).Build(w)

	return Setup{
		w,
		t,
	}
}

func (setup Setup) expectGroups(entity ecs.EntityID, expectedGroups groups.GroupsComponent) {
	groups, _ := setup.Groups().Component().Get(entity)
	if groups != expectedGroups {
		setup.T.Errorf("expected pos %v but has %v", expectedGroups, groups)
	}
}

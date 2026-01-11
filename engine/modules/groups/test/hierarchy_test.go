package test

import (
	"engine/modules/groups"
	"testing"
)

const (
	_ groups.Group = 1 << iota
	G1
)

func TestHierarchy(t *testing.T) {
	setup := NewSetup(t)
	defaultGroups := groups.EmptyGroups().Ptr().Enable(G1).Val()

	parent := setup.world.NewEntity()
	setup.groups.Component().Set(parent, defaultGroups)

	child := setup.world.NewEntity()
	setup.hierarchy.SetParent(child, parent)

	grandChild := setup.world.NewEntity()
	setup.hierarchy.SetParent(grandChild, child)

	setup.groups.Inherit().Set(grandChild, groups.InheritGroupsComponent{})
	setup.groups.Inherit().Set(child, groups.InheritGroupsComponent{})

	setup.expectGroups(child, defaultGroups)
	setup.expectGroups(grandChild, defaultGroups)
}

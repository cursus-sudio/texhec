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

	parent := setup.World.NewEntity()
	setup.Groups.Set(parent, defaultGroups)

	child := setup.World.NewEntity()
	setup.Hierarchy.SetParent(child, parent)

	grandChild := setup.World.NewEntity()
	setup.Hierarchy.SetParent(grandChild, child)

	setup.InheritGroups.Set(grandChild, groups.InheritGroupsComponent{})
	setup.InheritGroups.Set(child, groups.InheritGroupsComponent{})

	setup.expectGroups(child, defaultGroups)
	setup.expectGroups(grandChild, defaultGroups)
}

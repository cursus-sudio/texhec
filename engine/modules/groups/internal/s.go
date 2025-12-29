package internal

import (
	"engine/modules/groups"
	"engine/services/ecs"
)

func (t tool) calculateGroup(entity ecs.EntityID) (groups.GroupsComponent, bool) {
	def := groups.GroupsComponent{}
	parent, ok := t.world.Hierarchy().Parent(entity)
	if !ok {
		return def, false
	}
	groups, ok := t.groupsArray.Get(parent)
	if !ok {
		return def, false
	}
	return groups, ok
}

type save struct {
	entity ecs.EntityID
	groups groups.GroupsComponent
}

func (s tool) Init() {
	s.groupsArray.SetEmpty(groups.DefaultGroups())

	dirtySet := ecs.NewDirtySet()
	s.groupsArray.AddDependency(s.inheritArray)
	s.groupsArray.AddDependency(s.world.Hierarchy().Component())

	s.groupsArray.AddDirtySet(dirtySet)

	s.groupsArray.BeforeGet(func() {
		entities := dirtySet.Get()
		if len(entities) == 0 {
			return
		}
		children := []ecs.EntityID{}

		saves := []save{}

		for len(entities) != 0 || len(children) != 0 {
			if len(entities) == 0 {
				entities = children
				for _, save := range saves {
					s.groupsArray.Set(save.entity, save.groups)
				}

				dirtySet.Clear()
				children = nil
				saves = nil
			}
			entity := entities[0]
			entities = entities[1:]

			groups, ok := s.calculateGroup(entity)
			if !ok {
				continue
			}
			if originalGroups, ok := s.groupsArray.Get(entity); ok && groups == originalGroups {
				continue
			}
			saves = append(saves, save{
				entity: entity,
				groups: groups,
			})

			for _, child := range s.world.Hierarchy().Children(entity).GetIndices() {
				if _, ok := s.inheritArray.Get(child); ok {
					children = append(children, child)
				}
			}
		}

		for _, save := range saves {
			s.groupsArray.Set(save.entity, save.groups)
		}
		dirtySet.Clear()
	})
}

package internal

import (
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/services/ecs"
	"engine/services/logger"
)

type s struct {
	logger logger.Logger

	world     ecs.World
	hierarchy hierarchy.Interface

	hierarchyArray ecs.ComponentsArray[hierarchy.Component]
	inheritArray   ecs.ComponentsArray[groups.InheritGroupsComponent]
	groupsArray    ecs.ComponentsArray[groups.GroupsComponent]
}

func NewSystem(
	logger logger.Logger,
	parentToolFactory ecs.ToolFactory[hierarchy.HierarchyTool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := s{
			logger,
			w,
			parentToolFactory.Build(w).Hierarchy(),
			ecs.GetComponentsArray[hierarchy.Component](w),
			ecs.GetComponentsArray[groups.InheritGroupsComponent](w),
			ecs.GetComponentsArray[groups.GroupsComponent](w),
		}
		return s.Init()
	})
}

func (s s) calculateGroup(entity ecs.EntityID) (groups.GroupsComponent, bool) {
	def := groups.GroupsComponent{}
	parent, ok := s.hierarchy.Parent(entity)
	if !ok {
		return def, false
	}
	groups, ok := s.groupsArray.GetComponent(parent)
	if !ok {
		return def, false
	}
	return groups, ok
}

type save struct {
	entity ecs.EntityID
	groups groups.GroupsComponent
}

func (s s) Init() error {
	dirtySet := ecs.NewDirtySet()
	s.groupsArray.AddDependency(s.inheritArray)
	s.groupsArray.AddDependency(s.hierarchyArray)

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
					s.groupsArray.SaveComponent(save.entity, save.groups)
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
			if originalGroups, ok := s.groupsArray.GetComponent(entity); ok && groups == originalGroups {
				continue
			}
			saves = append(saves, save{
				entity: entity,
				groups: groups,
			})

			for _, child := range s.hierarchy.Children(entity).GetIndices() {
				if _, ok := s.inheritArray.GetComponent(child); ok {
					children = append(children, child)
				}
			}
		}

		for _, save := range saves {
			s.groupsArray.SaveComponent(save.entity, save.groups)
		}
		dirtySet.Clear()
	})

	return nil
}

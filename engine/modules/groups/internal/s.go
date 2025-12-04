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
	hierarchy hierarchy.Transaction

	inheritArray ecs.ComponentsArray[groups.InheritGroupsComponent]
	groupsArray  ecs.ComponentsArray[groups.GroupsComponent]

	groupsTransaction ecs.ComponentsArrayTransaction[groups.GroupsComponent]
}

func NewSystem(
	logger logger.Logger,
	parentToolFactory ecs.ToolFactory[hierarchy.Tool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		inheritArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](w)
		groupsArray := ecs.GetComponentsArray[groups.GroupsComponent](w)
		s := s{
			logger,
			w,
			parentToolFactory.Build(w).Transaction(),
			inheritArray,
			groupsArray,
			groupsArray.Transaction(),
		}
		return s.Init()
	})
}

func (s s) Init() error {
	onParentUpsert := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			groups, err := s.groupsArray.GetComponent(entity)
			if err != nil {
				continue
			}
			parentObject := s.hierarchy.GetObject(entity)
			children := parentObject.Children()
			for _, child := range children.GetIndices() {
				_, err := s.inheritArray.GetComponent(child)
				if err != nil {
					continue
				}
				s.groupsTransaction.SaveComponent(child, groups)
			}
		}
		s.logger.Warn(ecs.FlushMany(s.groupsTransaction))
	}
	s.groupsArray.OnAdd(onParentUpsert)
	s.groupsArray.OnChange(onParentUpsert)

	onChildUpsert := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			parentObject := s.hierarchy.GetObject(entity)
			parent, err := parentObject.Parent().Get()
			if err != nil {
				continue
			}
			parentGroup, err := s.groupsArray.GetComponent(parent.Parent)
			if err != nil {
				continue
			}
			s.groupsTransaction.SaveComponent(entity, parentGroup)
		}
		ecs.FlushMany(s.groupsTransaction)
	}
	childQuery := s.world.Query().
		Track(hierarchy.ParentComponent{}).
		Require(groups.InheritGroupsComponent{}).
		Build()
	childQuery.OnAdd(onChildUpsert)
	childQuery.OnChange(onChildUpsert)

	return nil
}

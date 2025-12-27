package hierarchytool

import (
	"engine/modules/hierarchy"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

type parentComponent struct{}

type tool struct {
	logger logger.Logger

	world          ecs.World
	hierarchyArray ecs.ComponentsArray[hierarchy.Component]
	parentArray    ecs.ComponentsArray[parentComponent]

	dirtySet ecs.DirtySet

	parents      datastructures.SparseArray[ecs.EntityID, ecs.EntityID]
	children     datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
	flatChildren datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
}

func NewTool(logger logger.Logger) hierarchy.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w hierarchy.World) hierarchy.HierarchyTool {
		mutex.Lock()
		defer mutex.Unlock()
		if tool, ok := ecs.GetGlobal[tool](w); ok {
			return tool
		}

		t := tool{
			logger,
			w,
			ecs.GetComponentsArray[hierarchy.Component](w),
			ecs.GetComponentsArray[parentComponent](w),
			ecs.NewDirtySet(),
			datastructures.NewSparseArray[ecs.EntityID, ecs.EntityID](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
		}
		w.SaveGlobal(t)
		t.hierarchyArray.AddDependency(t.parentArray)
		t.hierarchyArray.AddDirtySet(t.dirtySet)

		return t
	})
}

func (t tool) Hierarchy() hierarchy.Interface {
	return t
}
func (t tool) Component() ecs.ComponentsArray[hierarchy.Component] {
	return t.hierarchyArray
}

func (t tool) IsChildOf(child ecs.EntityID, wantedParent ecs.EntityID) bool {
	for {
		parent, ok := t.Parent(child)
		if !ok {
			return false
		}
		child = parent
		if child == wantedParent {
			return true
		}
	}
}

func (t tool) SetParent(child ecs.EntityID, parent ecs.EntityID) {
	t.hierarchyArray.Set(child, hierarchy.NewParent(parent))
}

func (t tool) Parent(child ecs.EntityID) (ecs.EntityID, bool) {
	comp, ok := t.hierarchyArray.Get(child)
	return comp.Parent, ok
}

//

func (t tool) GetParents(child ecs.EntityID) datastructures.SparseSet[ecs.EntityID] {
	orderedParents := t.GetOrderedParents(child)

	parents := datastructures.NewSparseSet[ecs.EntityID]()
	for _, parent := range orderedParents {
		parents.Add(parent)
	}
	return parents
}

func (t tool) GetOrderedParents(child ecs.EntityID) []ecs.EntityID {
	parents := []ecs.EntityID{child}
	for {
		parent, ok := t.hierarchyArray.Get(child)
		if !ok {
			return parents[1:]
		}
		parents = append(parents, parent.Parent)
		if len(parents) != 1 && parents[0] == parent.Parent {
			return nil
		}
		child = parent.Parent
	}
}

//

func (t tool) SetChildren(parent ecs.EntityID, children ...ecs.EntityID) {
	t.BeforeGet()
	previousChildren := t.Children(parent).GetIndices()
	i := 0
	for _, child := range previousChildren {
		if child == children[i] {
			i++
			continue
		}
		t.hierarchyArray.Remove(child)
	}

	for i := i; i < len(children); i++ {
		t.SetParent(children[i], parent)
	}
}

//

func (t tool) Children(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	t.BeforeGet()
	children, ok := t.children.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return children
}

func (t tool) FlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	t.BeforeGet()
	flatChildren, ok := t.flatChildren.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return flatChildren
}

//

func (t tool) BeforeGet() {
	dirtyEntities := t.dirtySet.Get()
	if len(dirtyEntities) == 0 {
		return
	}

	for _, child := range dirtyEntities {
		t.handleEntityChange(child)
	}
	t.dirtySet.Clear()
}

func (t tool) handleEntityChange(entity ecs.EntityID) {
	parent, parentOk := t.parents.Get(entity)
	hierarchy, hierarchyOk := t.hierarchyArray.Get(entity)
	_, isParent := t.parentArray.Get(entity)

	// children added to parents as flat children
	inheritedChildren, ok := t.flatChildren.Get(entity)
	if !ok {
		inheritedChildren = datastructures.NewSparseSet[ecs.EntityID]()
	}
	entityFlatChildren := inheritedChildren.GetIndices()
	inheritedChildren.Add(entity)
	inheritedChildrenIndices := inheritedChildren.GetIndices()

	if !hierarchyOk || parent != hierarchy.Parent {
		// remove parent
		t.parents.Remove(entity)

		// remove as a child
		children, ok := t.children.Get(parent)
		if !ok { // this shouldn't occur and means invalid internal state
			goto skipRemovalInParents
		}
		children.Remove(entity)
		if len(children.GetIndices()) == 0 {
			t.parents.Remove(entity)
			t.children.Remove(parent)
			t.flatChildren.Remove(parent)
		}

		// remove as a grand parent
		parents := t.GetOrderedParents(entity)
		for _, parent := range parents {
			flatChildren, ok := t.flatChildren.Get(parent)
			if !ok {
				continue
			}
			for _, child := range inheritedChildrenIndices {
				flatChildren.Remove(child)
			}
			// parent flat children are already removed if would be empty and
			// other grand parents have parent as a child so they won't be empty
			// so we do not have to check are flat children empty
		}
	}

skipRemovalInParents:
	if !parentOk && parent != hierarchy.Parent {
		// add parent
		t.parents.Set(entity, hierarchy.Parent)

		// add as parent
		parentChildren, ok := t.children.Get(hierarchy.Parent)
		if !ok {
			// mark as parent
			t.parentArray.Set(hierarchy.Parent, parentComponent{})

			// add children
			parentChildren = datastructures.NewSparseSet[ecs.EntityID]()
			t.children.Set(hierarchy.Parent, parentChildren)
		}
		parentChildren.Add(entity)
		t.parents.Set(entity, parent)

		// add as grand child
		parents := t.GetOrderedParents(entity)
		for _, parent := range parents {
			children, ok := t.flatChildren.Get(parent)
			if !ok {
				children = datastructures.NewSparseSet[ecs.EntityID]()
				t.flatChildren.Set(parent, children)
			}
			for _, child := range inheritedChildrenIndices {
				children.Add(child)
			}
		}
	}
	// skipChildrenAdditionInParents:

	if !isParent {
		t.children.Remove(entity)
		t.flatChildren.Remove(entity)
		t.parents.Remove(entity)
		for _, child := range entityFlatChildren {
			t.world.RemoveEntity(child)
			t.flatChildren.Remove(child)
			t.children.Remove(child)
			t.parents.Remove(child)
		}
	}
	// skipChildrenRemoval:
}

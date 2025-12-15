package hierarchytool

import (
	"engine/modules/hierarchy"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

type tool struct {
	logger logger.Logger

	world          ecs.World
	hierarchyArray ecs.ComponentsArray[hierarchy.Component]

	dirtySet ecs.DirtySet

	childrenParent     datastructures.SparseArray[ecs.EntityID, ecs.EntityID]
	parentChildren     datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
	parentFlatChildren datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
}

func NewTool(logger logger.Logger) ecs.ToolFactory[hierarchy.Hierarchy] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) hierarchy.Hierarchy {
		mutex.Lock()
		defer mutex.Unlock()
		if tool, ok := ecs.GetGlobal[tool](w); ok {
			return tool
		}

		dirtySet := ecs.NewDirtySet()

		tool := tool{
			logger,
			w,
			ecs.GetComponentsArray[hierarchy.Component](w),
			dirtySet,
			datastructures.NewSparseArray[ecs.EntityID, ecs.EntityID](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
		}
		w.SaveGlobal(tool)
		tool.Init()

		return tool
	})
}

func (t tool) Hierarchy() hierarchy.Interface {
	return t
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
	t.hierarchyArray.SaveComponent(child, hierarchy.NewParent(parent))
}

func (t tool) Parent(child ecs.EntityID) (ecs.EntityID, bool) {
	comp, ok := t.hierarchyArray.GetComponent(child)
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
		parent, ok := t.hierarchyArray.GetComponent(child)
		if !ok {
			return parents[1:]
		}
		parents = append(parents, parent.Parent)
		if parents[0] == parent.Parent {
			return nil
		}
		child = parent.Parent
	}
}

//

func (t tool) Children(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	t.BeforeGet()
	children, ok := t.parentChildren.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return children
}

func (t tool) FlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	t.BeforeGet()
	flatChildren, ok := t.parentFlatChildren.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return flatChildren
}

//

func (t tool) Init() {
	t.hierarchyArray.AddDirtySet(t.dirtySet)
}

func (t tool) BeforeGet() {
	dirtyEntities := t.dirtySet.Get()

	for _, child := range dirtyEntities {
		comp, ok := t.hierarchyArray.GetComponent(child)
		if !ok {
			t.handleRemoval(child)
			continue
		}
		parent := comp.Parent

		// modify child parent
		t.childrenParent.Set(child, comp.Parent)

		// modify parent
		children, ok := t.parentChildren.Get(parent)
		if !ok {
			children = datastructures.NewSparseSet[ecs.EntityID]()
			t.parentChildren.Set(parent, children)
		}
		children.Add(child)

		// modify flat parents
		parents := t.GetOrderedParents(child)
		for _, parent := range parents {
			children, ok := t.parentFlatChildren.Get(parent)
			if !ok {
				children = datastructures.NewSparseSet[ecs.EntityID]()
				t.parentFlatChildren.Set(parent, children)
			}
			children.Add(child)
		}
	}
}

func (t tool) handleRemoval(entity ecs.EntityID) {
	flatChildren := []ecs.EntityID{entity}
	if flatChildrenSet, ok := t.parentFlatChildren.Get(entity); ok {
		flatChildren = append(flatChildren, flatChildrenSet.GetIndices()...)
	}

	// remove children
	if children, ok := t.parentChildren.Get(entity); ok {
		for _, child := range children.GetIndices() {
			t.world.RemoveEntity(child)
		}
		t.parentChildren.Remove(entity)
		t.parentFlatChildren.Remove(entity)
	}

	// remove as a child
	parent, ok := t.childrenParent.Get(entity)
	if ok {
		if children, ok := t.parentChildren.Get(parent); ok {
			children.Remove(entity)
		}
	}
	for ok {
		if children, ok := t.parentFlatChildren.Get(parent); ok {
			for _, child := range flatChildren {
				children.Remove(child)
			}
		}

		parent, ok = t.childrenParent.Get(parent)
	}
}

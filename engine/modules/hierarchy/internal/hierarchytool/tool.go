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

	world       ecs.World
	parentArray ecs.ComponentsArray[hierarchy.ParentComponent]
	mutex       *sync.Mutex

	parentChildren     datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
	parentFlatChildren datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
}

func NewTool(logger logger.Logger) ecs.ToolFactory[hierarchy.Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) hierarchy.Tool {
		mutex.Lock()
		if tool, err := ecs.GetGlobal[tool](w); err == nil {
			mutex.Unlock()
			return tool
		}

		tool := tool{
			logger,
			w,
			ecs.GetComponentsArray[hierarchy.ParentComponent](w),
			mutex,
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
		}
		w.SaveGlobal(tool)
		mutex.Unlock()

		tool.parentArray.OnAdd(tool.Upsert)
		tool.parentArray.OnChange(tool.Upsert)
		tool.parentArray.OnRemoveComponents(tool.Remove)

		return tool
	})
}

func (t tool) Transaction() hierarchy.Transaction {
	return newTransaction(t)
}

func (t tool) GetOrderedParents(comp hierarchy.ParentComponent) []ecs.EntityID {
	parents := []ecs.EntityID{comp.Parent}
	child := comp.Parent
	for {
		parent, err := t.parentArray.GetComponent(child)
		if err != nil {
			return parents
		}
		parents = append(parents, parent.Parent)
		child = parent.Parent
	}
}

func (t tool) Upsert(ei []ecs.EntityID) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, child := range ei {
		comp, err := t.parentArray.GetComponent(child)
		if err != nil {
			continue
		}

		// parent
		children, ok := t.parentChildren.Get(comp.Parent)
		if !ok {
			children = datastructures.NewSparseSet[ecs.EntityID]()
			t.parentChildren.Set(comp.Parent, children)
		}
		children.Add(child)

		// flat parent

		for _, parent := range t.GetOrderedParents(comp) {
			children, ok := t.parentFlatChildren.Get(parent)
			if !ok {
				children = datastructures.NewSparseSet[ecs.EntityID]()
				t.parentFlatChildren.Set(comp.Parent, children)
			}
			children.Add(child)
		}
	}
}

func (t tool) Remove(ei []ecs.EntityID, components []hierarchy.ParentComponent) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, child := range ei {
		// handle orphaned children
		if children, ok := t.parentChildren.Get(child); ok {
			for _, child := range children.GetIndices() {
				t.world.RemoveEntity(child)
			}
			t.parentChildren.Remove(child)
			t.parentFlatChildren.Remove(child)
		}

		// remove from flat parents array
		for _, parent := range t.GetOrderedParents(components[i]) {
			children, ok := t.parentFlatChildren.Get(parent)
			if !ok {
				continue
			}
			if removed := children.Remove(child); !removed {
				continue
			}
			if len(children.GetIndices()) == 0 {
				t.parentFlatChildren.Remove(parent)
			}
		}

		// remove from parent array
		parentComponent := components[i]
		parent := parentComponent.Parent

		if children, ok := t.parentChildren.Get(parent); ok {
			if removed := children.Remove(child); removed && len(children.GetIndices()) == 0 {
				t.parentChildren.Remove(parent)
			}
		}
	}
}

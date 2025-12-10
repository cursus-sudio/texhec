package hierarchytool

import (
	"engine/modules/hierarchy"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"errors"
	"fmt"
	"sync"
)

type tool struct {
	logger logger.Logger

	world       ecs.World
	parentArray ecs.ComponentsArray[hierarchy.ParentComponent]

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
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
			datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),
		}
		w.SaveGlobal(tool)
		mutex.Unlock()

		tool.parentArray.OnAdd(tool.Upsert)
		tool.parentArray.OnChange(tool.Upsert)
		tool.parentArray.BeforeRemove(tool.Remove)

		return tool
	})
}

func (t tool) Transaction() hierarchy.Transaction {
	return newTransaction(t)
}

func (t tool) GetOrderedParents(comp hierarchy.ParentComponent) ([]ecs.EntityID, error) {
	parents := []ecs.EntityID{comp.Parent}
	child := comp.Parent
	for {
		parent, err := t.parentArray.GetComponent(child)
		if err != nil {
			return parents, nil
		}
		parents = append(parents, parent.Parent)
		if parents[0] == parent.Parent {
			return nil, errors.Join(
				hierarchy.ErrParentCycle,
				fmt.Errorf("cycle is %v", parents),
			)
		}
		child = parent.Parent
	}
}

func (t tool) Upsert(ei []ecs.EntityID) {
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

		parents, err := t.GetOrderedParents(comp)
		if err != nil {
			t.logger.Warn(err)
			continue
		}

		for _, parent := range parents {
			children, ok := t.parentFlatChildren.Get(parent)
			if !ok {
				children = datastructures.NewSparseSet[ecs.EntityID]()
				t.parentFlatChildren.Set(comp.Parent, children)
			}
			children.Add(child)
		}
	}
}

func (t tool) Remove(ei []ecs.EntityID) {
	for _, child := range ei {
		parentComponent, err := t.parentArray.GetComponent(child)
		if err != nil {
			t.logger.Warn(err)
			continue
		}
		// handle orphaned children
		if children, ok := t.parentChildren.Get(child); ok {
			for _, child := range children.GetIndices() {
				t.world.RemoveEntity(child)
			}
			t.parentChildren.Remove(child)
			t.parentFlatChildren.Remove(child)
		}

		// remove from flat parents array
		parents, err := t.GetOrderedParents(parentComponent)
		if err != nil {
			t.logger.Warn(err)
			continue
		}
		for _, parent := range parents {
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
		parent := parentComponent.Parent

		if children, ok := t.parentChildren.Get(parent); ok {
			if removed := children.Remove(child); removed && len(children.GetIndices()) == 0 {
				t.parentChildren.Remove(parent)
			}
		}
	}
}

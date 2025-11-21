package manytoone

import (
	"frontend/modules/relation"
	"shared/services/datastructures"
	"shared/services/ecs"
	"sync"
)

type manyToOne[Component any] struct {
	world          ecs.World
	childrenArray  ecs.ComponentsArray[Component]
	mutex          *sync.Mutex
	parentChildren datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]

	upsertListeners []func([]ecs.EntityID)
	removeListeners []func([]ecs.EntityID)

	getParent func(Component) ecs.EntityID
}

func newIndex[Component any](
	w ecs.World,
	componentParent func(Component) ecs.EntityID,
) relation.EntityToEntitiesTool[Component] {
	childrenArray := ecs.GetComponentsArray[Component](w)

	indexGlobal := manyToOne[Component]{
		world:          w,
		childrenArray:  childrenArray,
		mutex:          &sync.Mutex{},
		parentChildren: datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]](),

		upsertListeners: make([]func([]ecs.EntityID), 0),
		removeListeners: make([]func([]ecs.EntityID), 0),

		getParent: componentParent,
	}
	w.SaveGlobal(indexGlobal)

	childrenArray.OnAdd(indexGlobal.Upsert)
	childrenArray.OnChange(indexGlobal.Upsert)
	childrenArray.OnRemoveComponents(indexGlobal.Remove)

	return indexGlobal
}
func (i manyToOne[IndexType]) GetMany(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	children, ok := i.parentChildren.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return children
}

func (i manyToOne[IndexType]) OnUpsert(listener func([]ecs.EntityID)) {
	if values := i.parentChildren.GetIndices(); len(values) != 0 {
		listener(values)
	}
	i.upsertListeners = append(i.upsertListeners, listener)
	i.world.SaveGlobal(i)
}

func (i manyToOne[IndexType]) OnRemove(listener func([]ecs.EntityID)) {
	i.removeListeners = append(i.removeListeners, listener)
	i.world.SaveGlobal(i)
}

func (i manyToOne[Component]) Upsert(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	upsertedParents := datastructures.NewSparseSet[ecs.EntityID]()
	for _, child := range ei {
		comp, err := i.childrenArray.GetComponent(child)
		if err != nil {
			continue
		}
		parent := i.getParent(comp)
		children, ok := i.parentChildren.Get(parent)
		if !ok {
			children = datastructures.NewSparseSet[ecs.EntityID]()
			i.parentChildren.Set(parent, children)
		}
		if addedChild := children.Add(child); addedChild {
			upsertedParents.Add(parent)
		}
	}
	if entities := upsertedParents.GetIndices(); len(entities) != 0 {
		for _, listener := range i.upsertListeners {
			listener(entities)
		}
	}
}

func (r manyToOne[Components]) Remove(ei []ecs.EntityID, components []Components) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	changedParents := datastructures.NewSparseSet[ecs.EntityID]()
	removedParents := datastructures.NewSparseSet[ecs.EntityID]()
	for i, child := range ei {
		parent := r.getParent(components[i])
		children, ok := r.parentChildren.Get(parent)
		if !ok {
			continue
		}
		if removed := children.Remove(child); !removed {
			continue
		}
		if len(children.GetIndices()) == 0 {
			r.parentChildren.Remove(parent)
			removedParents.Add(parent)
		} else {
			changedParents.Add(parent)
		}
	}
	if entities := changedParents.GetIndices(); len(entities) != 0 {
		for _, listener := range r.upsertListeners {
			listener(entities)
		}
	}
	if entities := removedParents.GetIndices(); len(entities) != 0 {
		for _, listener := range r.removeListeners {
			listener(entities)
		}
	}
}

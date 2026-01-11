package hierarchyservice

import (
	"engine/modules/hierarchy"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type parentComponent struct{}

type service struct {
	Logger logger.Logger `inject:"1"`

	World          ecs.World `inject:"1"`
	hierarchyArray ecs.ComponentsArray[hierarchy.Component]
	parentArray    ecs.ComponentsArray[parentComponent]

	parents      datastructures.SparseArray[ecs.EntityID, ecs.EntityID]
	children     datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
	flatChildren datastructures.SparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]
}

func NewService(c ioc.Dic) hierarchy.Service {
	t := ioc.GetServices[*service](c)

	t.hierarchyArray = ecs.GetComponentsArray[hierarchy.Component](t.World)
	t.parentArray = ecs.GetComponentsArray[parentComponent](t.World)
	t.parents = datastructures.NewSparseArray[ecs.EntityID, ecs.EntityID]()
	t.children = datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]()
	t.flatChildren = datastructures.NewSparseArray[ecs.EntityID, datastructures.SparseSet[ecs.EntityID]]()

	t.hierarchyArray.OnUpsert(t.handleHierarchyChange)
	t.hierarchyArray.OnRemove(t.handleHierarchyChange)
	t.parentArray.OnRemove(t.handleParentChange)

	return t
}

func (t *service) Component() ecs.ComponentsArray[hierarchy.Component] {
	return t.hierarchyArray
}

func (t *service) IsChildOf(child ecs.EntityID, wantedParent ecs.EntityID) bool {
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

func (t *service) SetParent(child ecs.EntityID, parent ecs.EntityID) {
	t.hierarchyArray.Set(child, hierarchy.NewParent(parent))
}

func (t *service) Parent(child ecs.EntityID) (ecs.EntityID, bool) {
	comp, ok := t.hierarchyArray.Get(child)
	return comp.Parent, ok
}

//

func (t *service) GetParents(child ecs.EntityID) datastructures.SparseSet[ecs.EntityID] {
	orderedParents := t.GetOrderedParents(child)

	parents := datastructures.NewSparseSet[ecs.EntityID]()
	for _, parent := range orderedParents {
		parents.Add(parent)
	}
	return parents
}

func (t *service) GetOrderedParents(child ecs.EntityID) []ecs.EntityID {
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

func (t *service) GetOrderedPreviousParents(child ecs.EntityID) []ecs.EntityID {
	parents := []ecs.EntityID{child}
	for {
		parent, ok := t.parents.Get(child)
		if !ok {
			return parents[1:]
		}
		parents = append(parents, parent)
		if len(parents) != 1 && parents[0] == parent {
			return nil
		}
		child = parent
	}
}

//

func (t *service) SetChildren(parent ecs.EntityID, children ...ecs.EntityID) {
	previousChildren := t.Children(parent).GetIndices()
	for _, child := range previousChildren {
		t.hierarchyArray.Remove(child)
	}

	for i := 0; i < len(children); i++ {
		t.SetParent(children[i], parent)
	}
}

//

func (t *service) Children(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	children, ok := t.children.Get(parent)
	if !ok {
		return datastructures.NewSparseSet[ecs.EntityID]()
	}
	return children
}

func (t *service) GetFlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	if flatChildren, ok := t.flatChildren.Get(parent); ok {
		return flatChildren
	}
	flatChildren := datastructures.NewSparseSet[ecs.EntityID]()

	children, ok := t.children.Get(parent)
	if !ok {
		return flatChildren
	}

	childrens := []datastructures.SparseSet[ecs.EntityID]{children}

	for len(childrens) != 0 {
		children := childrens[0]
		childrens = childrens[1:]

		for _, child := range children.GetIndices() {
			if added := flatChildren.Add(child); !added {
				continue
			}
			children, ok := t.children.Get(child)
			if !ok {
				continue
			}
			childrens = append(childrens, children)
		}
	}

	t.flatChildren.Set(parent, flatChildren)
	return flatChildren
}

func (t *service) FlatChildren(parent ecs.EntityID) datastructures.SparseSetReader[ecs.EntityID] {
	return t.GetFlatChildren(parent)
}

//

func (t *service) handleHierarchyChange(child ecs.EntityID) {
	previousParent, previousParentOk := t.parents.Get(child)
	hierarchy, nextParentOk := t.hierarchyArray.Get(child)
	if previousParentOk == nextParentOk && hierarchy.Parent == previousParent {
		return
	}

	if previousParentOk { // remove in parents
		t.parents.Remove(child)

		for _, parent := range t.GetOrderedPreviousParents(child) {
			t.flatChildren.Remove(parent)
		}

		// remove as a child
		children, ok := t.children.Get(previousParent)
		if !ok { // this shouldn't occur and means invalid internal state
			goto addCurrentParent
		}
		children.Remove(child)
		if len(children.GetIndices()) == 0 {
			t.children.Remove(previousParent)
		}
	}

addCurrentParent:
	nextParent := hierarchy.Parent
	if nextParentOk { // add in parents
		// add parent
		t.parents.Set(child, nextParent)

		// add as parent
		parentChildren, ok := t.children.Get(nextParent)
		if !ok {
			// mark as parent
			t.parentArray.Set(nextParent, parentComponent{})

			// add children
			parentChildren = datastructures.NewSparseSet[ecs.EntityID]()
			t.children.Set(nextParent, parentChildren)
		}
		parentChildren.Add(child)
	}
}

func (t *service) handleParentChange(parent ecs.EntityID) {
	if _, isParent := t.parentArray.Get(parent); isParent {
		return
	}

	children := t.GetFlatChildren(parent)

	for _, parent := range t.GetOrderedParents(parent) {
		t.flatChildren.Remove(parent)
	}

	t.children.Remove(parent)
	t.flatChildren.Remove(parent)
	for _, child := range children.GetIndices() {
		t.flatChildren.Remove(child)
		t.children.Remove(child)
		t.parents.Remove(child)
		t.World.RemoveEntity(child)
	}
}

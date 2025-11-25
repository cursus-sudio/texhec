package ecs

import (
	"engine/services/datastructures"
)

// interface

type LiveQuery interface {
	// on add listener add all entities are passed to it
	OnAdd(func([]EntityID))
	OnChange(func([]EntityID))
	OnRemove(func([]EntityID))

	Entities() []EntityID
}

// impl

type liveQuery struct {
	entities datastructures.SparseSet[EntityID]
	onRemove []func([]EntityID)
	onChange []func([]EntityID)
	onAdd    []func([]EntityID)
}

func newLiveQuery(
	componentsImpl *componentsImpl,
	required []ComponentType,
	tracked []ComponentType,
	forbidden []ComponentType,
) *liveQuery {
	liveQuery := &liveQuery{
		entities: datastructures.NewSparseSet[EntityID](),
	}

	requiredArrays := make([]arraysSharedInterface, 0, len(required))
	forbiddenArrays := make([]arraysSharedInterface, 0, len(forbidden))

	tryAddEntities := func(ei []EntityID) {
		// this optimizes in case when not all required arrays aren't there
		if len(requiredArrays) != len(required) {
			return
		}
		addedEntities := []EntityID{}
	entityLoop:
		for _, entity := range ei {
			if ok := liveQuery.entities.Get(entity); ok {
				continue
			}

			// check requirements
			for _, arr := range requiredArrays {
				if _, err := arr.GetAnyComponent(entity); err != nil {
					continue entityLoop
				}
			}
			for _, arr := range forbiddenArrays {
				if _, err := arr.GetAnyComponent(entity); err == nil {
					continue entityLoop
				}
			}

			if added := liveQuery.entities.Add(entity); added {
				addedEntities = append(addedEntities, entity)
			}
		}

		if len(addedEntities) == 0 {
			return
		}
		for _, listener := range liveQuery.onAdd {
			listener(addedEntities)
		}
	}
	changeEntities := func(ei []EntityID) {
		changedEntities := []EntityID{}
		for _, entity := range ei {
			ok := liveQuery.entities.Get(entity)
			if ok {
				changedEntities = append(changedEntities, entity)
			}
		}

		if len(changedEntities) == 0 {
			return
		}
		for _, listener := range liveQuery.onChange {
			listener(changedEntities)
		}
	}
	removeEntities := func(ei []EntityID) {
		removedEntities := []EntityID{}
		for _, entity := range ei {
			removed := liveQuery.entities.Remove(entity)
			if removed {
				removedEntities = append(removedEntities, entity)
			}
		}

		if len(removedEntities) == 0 {
			return
		}
		for _, listener := range liveQuery.onRemove {
			listener(removedEntities)
		}
	}

	for _, required := range required {
		componentsImpl.storage.whenArrExists(required, func(arr arraysSharedInterface) {
			requiredArrays = append(requiredArrays, arr)
			arr.OnAdd(tryAddEntities)
			arr.OnChange(changeEntities)
			arr.OnRemove(removeEntities)
		})
	}
	for _, tracked := range tracked {
		componentsImpl.storage.whenArrExists(tracked, func(arr arraysSharedInterface) {
			arr.OnAdd(changeEntities)
			arr.OnChange(changeEntities)
			arr.OnRemove(changeEntities)
		})
	}
	for _, forbid := range forbidden {
		componentsImpl.storage.whenArrExists(forbid, func(arr arraysSharedInterface) {
			forbiddenArrays = append(forbiddenArrays, arr)
			arr.OnAdd(removeEntities)
			arr.OnChange(changeEntities)
			arr.OnRemove(tryAddEntities)
		})
	}
	return liveQuery
}

func (query *liveQuery) OnAdd(listener func([]EntityID)) {
	entities := query.Entities()
	if len(entities) != 0 {
		listener(entities)
	}
	query.onAdd = append(query.onAdd, listener)
}

func (query *liveQuery) OnChange(listener func([]EntityID)) {
	query.onChange = append(query.onChange, listener)
}

func (query *liveQuery) OnRemove(listener func([]EntityID)) {
	query.onRemove = append(query.onRemove, listener)
}

func (query *liveQuery) Entities() []EntityID { return query.entities.GetIndices() }

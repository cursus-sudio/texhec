package ecs

import (
	"frontend/services/datastructures"
	"sort"
	"strings"
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

type queryKey string

func newQueryKey(components []ComponentType) queryKey {
	resultLen := 0
	elements := make([]string, len(components))
	for i, component := range components {
		element := component.componentType.String()
		resultLen += len(element) + 1
		elements[i] = element
	}
	sort.Strings(elements)
	builder := strings.Builder{}
	builder.Grow(resultLen)
	for _, element := range elements {
		builder.WriteString(element)
		builder.WriteString(",")
	}
	return queryKey(builder.String())
}

//

type liveQuery struct {
	dependencies datastructures.Set[ComponentType]
	entities     datastructures.SparseSet[EntityID]
	onRemove     []func([]EntityID)
	onChange     []func([]EntityID)
	onAdd        []func([]EntityID)
}

func newLiveQuery(componentTypes []ComponentType) *liveQuery {
	dependencies := datastructures.NewSet[ComponentType]()
	for _, componentType := range componentTypes {
		dependencies.Add(componentType)
	}
	return &liveQuery{
		dependencies: dependencies,
		entities:     datastructures.NewSparseSet[EntityID](),
	}
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

//

func (query *liveQuery) AddedEntities(entities []EntityID) {
	added := make([]EntityID, 0, len(entities))
	for _, entity := range entities {
		if added := query.entities.Add(entity); !added {
			continue
		}
		added = append(added, entity)
	}
	if len(added) != 0 {
		for _, listener := range query.onAdd {
			listener(entities)
		}
	}
}

func (query *liveQuery) Changed(entities []EntityID) {
	for _, listener := range query.onChange {
		listener(entities)
	}
}

func (query *liveQuery) RemovedEntities(entities []EntityID) {
	removed := make([]EntityID, 0, len(entities))
	for _, entity := range entities {
		if removed := query.entities.Remove(entity); !removed {
			continue
		}
		removed = append(removed, entity)
	}
	if len(removed) != 0 {
		for _, listener := range query.onRemove {
			listener(entities)
		}
	}
}

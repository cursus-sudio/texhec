package ecs

import (
	"frontend/services/datastructures"
	"sort"
	"strings"
)

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
	dependencies map[ComponentType]any // this is faster []ComponentType
	entities     datastructures.Set[EntityID]
	onRemove     []func([]EntityID)
	onChange     []func([]EntityID)
	onAdd        []func([]EntityID)
}

func newLiveQuery(
	componentTypes []ComponentType,
	res []EntityID,
) *liveQuery {
	dependencies := make(map[ComponentType]any, len(componentTypes))
	for _, componentType := range componentTypes {
		dependencies[componentType] = nil
	}
	entities := datastructures.NewSet[EntityID]()
	for _, entity := range res {
		entities.Add(entity)
	}
	return &liveQuery{
		dependencies: dependencies,
		entities:     entities,
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

func (query *liveQuery) Entities() []EntityID {
	return query.entities.Get()
}

//

func (query *liveQuery) RemoveEntity(entity EntityID) {
	index, ok := query.entities.GetIndex(entity)
	if !ok {
		return
	}
	query.entities.Remove(index)
	rmArgs := []EntityID{entity}
	for _, listener := range query.onRemove {
		listener(rmArgs)
	}
}

func (query *liveQuery) Changed(entities []EntityID) {
	for _, listener := range query.onChange {
		listener(entities)
	}
}

func (query *liveQuery) AddEntities(entities []EntityID) {
	query.entities.Add(entities...)
	for _, listener := range query.onAdd {
		listener(entities)
	}
}

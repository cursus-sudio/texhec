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
	dependencies   datastructures.Set[ComponentType] // this is faster []ComponentType
	entities       []uint32
	cachedEntities []EntityID
	onRemove       []func([]EntityID)
	onChange       []func([]EntityID)
	onAdd          []func([]EntityID)
}

func newLiveQuery(componentTypes []ComponentType) *liveQuery {
	dependencies := datastructures.NewSet[ComponentType]()
	for _, componentType := range componentTypes {
		dependencies.Add(componentType)
	}
	return &liveQuery{
		dependencies:   dependencies,
		entities:       []uint32{},
		cachedEntities: []EntityID{},
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
	return query.cachedEntities
}

//

func (query *liveQuery) AddedEntities(entities []EntityID) {
	for _, entity := range entities {
		entityIndex := entity.Index()
		for entityIndex >= len(query.entities) {
			query.entities = append(query.entities, noEntity)
		}
		query.entities[entityIndex] = uint32(len(query.cachedEntities))
		query.cachedEntities = append(query.cachedEntities, entity)
	}
	for _, listener := range query.onAdd {
		listener(entities)
	}
}

func (query *liveQuery) Changed(entities []EntityID) {
	for _, listener := range query.onChange {
		listener(entities)
	}
}

func (query *liveQuery) RemovedEntities(entities []EntityID) {
	for _, entity := range entities {
		entityIndex := entity.Index()
		cachedIndex := query.entities[entityIndex]

		query.entities[entityIndex] = noEntity

		if len(query.cachedEntities)-1 != int(cachedIndex) {
			movedComponentEntity := query.cachedEntities[len(query.cachedEntities)-1]
			query.cachedEntities[cachedIndex] = movedComponentEntity

			movedEntityIndex := movedComponentEntity.Index()
			query.entities[movedEntityIndex] = cachedIndex
		}

		query.cachedEntities = query.cachedEntities[:len(query.cachedEntities)-1]
	}
	for _, listener := range query.onRemove {
		listener(entities)
	}
}

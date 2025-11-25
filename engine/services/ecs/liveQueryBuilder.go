package ecs

import (
	"strings"
)

type LiveQueryBuilder interface {
	Require(...ComponentType) LiveQueryBuilder

	// is used to call onChange on liveQuery on componentType add, change or remove
	Track(...ComponentType) LiveQueryBuilder
	Forbid(...ComponentType) LiveQueryBuilder
	Build() LiveQuery
}

type liveQueryFactory struct {
	components *componentsImpl

	required  []ComponentType
	tracked   []ComponentType
	forbidden []ComponentType
}

func newLiveQueryFactory(components *componentsImpl) LiveQueryBuilder {
	return &liveQueryFactory{
		components: components,
	}
}

func (f *liveQueryFactory) Require(components ...ComponentType) LiveQueryBuilder {
	f.required = append(f.required, components...)
	return f
}

func (f *liveQueryFactory) Track(components ...ComponentType) LiveQueryBuilder {
	f.tracked = append(f.tracked, components...)
	return f
}

func (f *liveQueryFactory) Forbid(components ...ComponentType) LiveQueryBuilder {
	f.forbidden = append(f.forbidden, components...)
	return f
}

func (f *liveQueryFactory) Key() queryKey {
	s := strings.Join([]string{
		typesArrayTostring(f.required),
		typesArrayTostring(f.tracked),
		typesArrayTostring(f.forbidden),
	}, "|")

	return queryKey(s)
}

func (f *liveQueryFactory) Build() LiveQuery {
	key := f.Key()
	if query, ok := f.components.storage.cachedQueries[key]; ok {
		return query
	}
	query := newLiveQuery(f.components, f.required, f.tracked, f.forbidden)
	f.components.storage.cachedQueries[key] = query
	return query
}

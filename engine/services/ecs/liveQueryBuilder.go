package ecs

import (
	"strings"
)

type LiveQueryBuilder interface {
	Require(...any) LiveQueryBuilder

	// is used to call onChange on liveQuery on componentType add, change or remove
	Track(...any) LiveQueryBuilder
	Forbid(...any) LiveQueryBuilder
	Build() LiveQuery
}

type liveQueryFactory struct {
	components *componentsImpl

	required  []componentType
	tracked   []componentType
	forbidden []componentType
}

func newLiveQueryFactory(components *componentsImpl) LiveQueryBuilder {
	return &liveQueryFactory{
		components: components,
	}
}

func (f *liveQueryFactory) Require(components ...any) LiveQueryBuilder {
	for _, component := range components {
		f.required = append(f.required, getComponentType(component))
	}
	return f
}

func (f *liveQueryFactory) Track(components ...any) LiveQueryBuilder {
	for _, component := range components {
		f.tracked = append(f.tracked, getComponentType(component))
	}
	return f
}

func (f *liveQueryFactory) Forbid(components ...any) LiveQueryBuilder {
	for _, component := range components {
		f.forbidden = append(f.forbidden, getComponentType(component))
	}
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

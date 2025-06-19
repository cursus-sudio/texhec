package backendconnector

import (
	"backend/services/backendapi"
	"sync"
)

type Builder interface {
	DefaultConnection(defaultConnector func() backendapi.Backend)
	Build() Backend
}

type builder struct {
	defaultConnector func() backendapi.Backend
}

func NewBuilder() Builder {
	return &builder{}
}

func (builder *builder) DefaultConnection(defaultConnector func() backendapi.Backend) {
	builder.defaultConnector = defaultConnector
}

func (builder *builder) Build() Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           builder.defaultConnector(),
		getDefaultConnection: builder.defaultConnector,
	}
}

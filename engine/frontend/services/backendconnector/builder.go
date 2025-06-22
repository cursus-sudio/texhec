package backendconnector

import (
	"backend/services/api"
	"sync"
)

type Builder interface {
	DefaultConnection(defaultConnector func() api.Server)
	Build() Backend
}

type builder struct {
	defaultConnector func() api.Server
}

func NewBuilder() Builder {
	return &builder{}
}

func (builder *builder) DefaultConnection(defaultConnector func() api.Server) {
	builder.defaultConnector = defaultConnector
}

func (builder *builder) Build() Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           builder.defaultConnector(),
		getDefaultConnection: builder.defaultConnector,
	}
}

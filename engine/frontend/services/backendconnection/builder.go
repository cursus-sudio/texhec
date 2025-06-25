package backendconnection

import (
	"shared/utils/connection"
	"sync"
)

type Builder interface {
	DefaultConnection(defaultConnector func() connection.Connection)
	Build() Backend
}

type builder struct {
	defaultConnector func() connection.Connection
}

func NewBuilder() Builder {
	return &builder{}
}

func (builder *builder) DefaultConnection(defaultConnector func() connection.Connection) {
	builder.defaultConnector = defaultConnector
}

func (builder *builder) Build() Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           builder.defaultConnector(),
		getDefaultConnection: builder.defaultConnector,
	}
}

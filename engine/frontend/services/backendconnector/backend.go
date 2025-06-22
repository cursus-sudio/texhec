package backendconnector

import (
	"backend/services/api"
	"backend/utils/httperrors"
	"sync"
)

type Backend interface {
	Connect(api.Server) error
	Disconnect()
	Connection() api.Server
}

type backend struct {
	rwMutex              sync.RWMutex
	connection           api.Server
	getDefaultConnection func() api.Server
}

func NewBackend(getDefaultConnection func() api.Server) Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           getDefaultConnection(),
		getDefaultConnection: getDefaultConnection,
	}
}

func (backend *backend) Connect(connection api.Server) error {
	if connection == nil {
		return httperrors.Err400
	}
	backend.rwMutex.Lock()
	defer backend.rwMutex.Unlock()
	backend.connection = connection
	return nil
}

func (backend *backend) Disconnect() {
	backend.connection = backend.getDefaultConnection()
}

func (backend *backend) Connection() api.Server {
	if backend.connection == nil {

	}
	backend.rwMutex.RLock()
	defer backend.rwMutex.RUnlock()
	return backend.connection
}

package backendconnector

import (
	"backend/src/backendapi"
	"backend/src/utils/httperrors"
	"sync"
)

type Backend interface {
	Connect(backendapi.Backend) error
	Disconnect()
	Connection() backendapi.Backend
}

type backend struct {
	rwMutex              sync.RWMutex
	connection           backendapi.Backend
	getDefaultConnection func() backendapi.Backend
}

func NewBackend(getDefaultConnection func() backendapi.Backend) Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           getDefaultConnection(),
		getDefaultConnection: getDefaultConnection,
	}
}

func (backend *backend) Connect(connection backendapi.Backend) error {
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

func (backend *backend) Connection() backendapi.Backend {
	if backend.connection == nil {

	}
	backend.rwMutex.RLock()
	defer backend.rwMutex.RUnlock()
	return backend.connection
}

package backendconnection

import (
	"shared/utils/connection"
	"shared/utils/httperrors"
	"sync"
)

type Backend interface {
	Connect(connection.Connection) error
	Disconnect()
	Connection() connection.Connection
}

type backend struct {
	rwMutex              sync.RWMutex
	connection           connection.Connection
	getDefaultConnection func() connection.Connection
}

func NewBackend(getDefaultConnection func() connection.Connection) Backend {
	return &backend{
		rwMutex:              sync.RWMutex{},
		connection:           getDefaultConnection(),
		getDefaultConnection: getDefaultConnection,
	}
}

func (backend *backend) Connect(connection connection.Connection) error {
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

func (backend *backend) Connection() connection.Connection {
	if backend.connection == nil {

	}
	backend.rwMutex.RLock()
	defer backend.rwMutex.RUnlock()
	return backend.connection
}

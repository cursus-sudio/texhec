package clients

import (
	"backend/utils/httperrors"
	"sync"
)

type ClientID string

type Clients interface {
	Client(ClientID) (ClientEmitter, error)
	Connect(ClientID, Client) error
	Disconnect(ClientID)
}

type clients struct {
	clients     map[ClientID]Client
	clientMutex sync.Mutex
}

func NewClients() Clients {
	return &clients{
		clients:     map[ClientID]Client{},
		clientMutex: sync.Mutex{},
	}
}

func (clients *clients) Client(id ClientID) (ClientEmitter, error) {
	clients.clientMutex.Lock()
	defer clients.clientMutex.Unlock()
	client, ok := clients.clients[id]
	if ok {
		return client, nil
	}
	return nil, httperrors.Err404
}

func (clients *clients) Connect(id ClientID, client Client) error {
	clients.clientMutex.Lock()
	defer clients.clientMutex.Unlock()
	_, ok := clients.clients[id]
	if ok {
		return httperrors.Err409
	}
	clients.clients[id] = client
	return nil
}

func (clients *clients) Disconnect(id ClientID) {
	clients.clientMutex.Lock()
	defer clients.clientMutex.Unlock()
	client, ok := clients.clients[id]
	if !ok {
		return
	}
	client.Close()
	delete(clients.clients, id)

}

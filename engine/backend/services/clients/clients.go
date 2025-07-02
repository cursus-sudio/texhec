package clients

import (
	"shared/utils/httperrors"
	"sync"
)

type ConnectListener func(Client)

type ClientsBuilder interface {
	OnConnect(listener ConnectListener) ClientsBuilder
	Build() Clients
}

type clientBuilder struct {
	listeners []ConnectListener
}

func NewClientsBuilder() ClientsBuilder {
	return &clientBuilder{
		listeners: []ConnectListener{},
	}
}

func (b *clientBuilder) OnConnect(listener ConnectListener) ClientsBuilder {
	b.listeners = append(b.listeners, listener)
	return b
}

func (b *clientBuilder) Build() Clients {
	onConnectListeners := b.listeners
	return &clients{
		onConnect: func(c Client) {
			for _, listener := range onConnectListeners {
				listener(c)
			}
		},
		clients: map[ClientID]Client{},
	}
}

//

type Clients interface {
	AllClients() []Client

	// 404
	ClientById(ClientID) (Client, error)

	// 409
	Connect(Client) error
}

type clients struct {
	mutex     sync.Mutex
	onConnect func(Client)
	clients   map[ClientID]Client
}

func (clients *clients) AllClients() []Client {
	r := make([]Client, 0, len(clients.clients))
	for _, client := range clients.clients {
		r = append(r, client)
	}
	return r
}

func (clients *clients) ClientById(id ClientID) (Client, error) {
	client, ok := clients.clients[id]
	if !ok {
		var client Client
		return client, httperrors.Err404
	}
	return client, nil
}

func (clients *clients) Disconnect(id ClientID) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	delete(clients.clients, id)
}

func (clients *clients) Connect(client Client) error {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()
	if _, ok := clients.clients[client.ID]; ok {
		return httperrors.Err409
	}
	clients.clients[client.ID] = client
	client.Connection.OnClose(func() { clients.Disconnect(client.ID) })
	return nil
}

package clients

import "shared/utils/httperrors"

type SessionClient interface {
	// 404
	Client() (Client, error)

	// 409
	UseClient(Client) error
}

type sessionClient struct {
	clients Clients
	client  *Client
}

func NewSessionClient(clients Clients) SessionClient {
	return &sessionClient{
		clients: clients,
		client:  nil,
	}
}

func (session *sessionClient) Client() (Client, error) {
	if session.client == nil {
		var client Client
		return client, httperrors.Err404
	}
	return *session.client, nil
}

func (session *sessionClient) UseClient(client Client) error {
	if session.client != nil {
		return httperrors.Err409
	}
	heapClient := client
	session.client = &heapClient
	session.clients.Connect(client)
	return nil
}

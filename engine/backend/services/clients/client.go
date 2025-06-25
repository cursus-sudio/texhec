package clients

import "shared/utils/connection"

type ClientID string

type Client struct {
	ID         ClientID
	Connection connection.Connection
}

func NewClient(id ClientID, connection connection.Connection) Client {
	return Client{
		ID:         id,
		Connection: connection,
	}
}

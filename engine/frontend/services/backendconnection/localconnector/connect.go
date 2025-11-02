package localconnector

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Connector interface {
	Connect() connection.Connection
}

type connector struct {
	connect func() connection.Connection
}

func newConnector(connect func() connection.Connection) Connector {
	return &connector{
		connect: connect,
	}
}

func (connector *connector) Connect() connection.Connection {
	return connector.connect()
}

type pkg struct {
	connect func(connection.Connection) connection.Connection
}

func Package(
	connect func(connection.Connection) connection.Connection,
) ioc.Pkg {
	return pkg{connect: connect}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Connector {
		return newConnector(func() connection.Connection {
			return pkg.connect(ioc.Get[connection.Connection](c))
		})
	})
}

package tcp

import (
	"frontend/services/backendconnection"
	"net"
	"shared/services/api/netconnection"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	network string
}

func Package(
	network string,
) ioc.Pkg {
	return pkg{
		network: network,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Connect {
		return newConnect(
			pkg.network,
			ioc.Get[netconnection.NetConnection](c),
			ioc.Get[backendconnection.Backend](c),
		)
	})
	ioc.RegisterDependency[Connect, netconnection.NetConnection](b)
	ioc.RegisterDependency[Connect, backendconnection.Backend](b)
}

//

type Connect interface {
	Connect(address string) error
}

type connect struct {
	network string
	net     netconnection.NetConnection
	backend backendconnection.Backend
}

func newConnect(
	network string,
	net netconnection.NetConnection,
	backend backendconnection.Backend,
) Connect {
	return &connect{
		network: network,
		net:     net,
		backend: backend,
	}
}

func (s connect) Connect(address string) error {
	conn, err := net.Dial(s.network, address)
	if err != nil {
		return err
	}
	var connection connection.Connection
	connection = s.net.Connect(conn, func() {
		if s.backend.Connection() == connection {
			s.backend.Disconnect()
		}
	})
	s.backend.Connect(connection)
	return nil
}

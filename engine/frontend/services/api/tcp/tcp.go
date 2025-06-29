package tcp

import (
	"frontend/services/backendconnection"
	"net"
	"shared/services/connpkg/netconnection"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	network string
}

func Package(
	network string,
) Pkg {
	return Pkg{
		network: network,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Connect {
		return newConnect(
			pkg.network,
			ioc.Get[netconnection.NetConnection](c),
			ioc.Get[backendconnection.Backend](c),
		)
	})
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

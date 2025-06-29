package tcp

import (
	"backend/services/clients"
	"backend/services/logger"
	"backend/services/scopes"
	"fmt"
	"net"
	"shared/services/api/netconnection"
	"shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
)

func (pkg Pkg) hostServer(c ioc.Dic) {
	logger := ioc.Get[logger.Logger](c)
	runtime := ioc.Get[runtime.Runtime](c)
	l, err := net.Listen(pkg.connType, pkg.host+":"+pkg.port)
	if err != nil {
		err := fmt.Errorf("Error listening: %s", err.Error())
		logger.Error(err)
		runtime.Stop()
		return
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Printf("%s\n", logger)
	logger.Info(
		fmt.Sprintf("Listening on %s:%s",
			pkg.host,
			pkg.port,
		),
	)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		c := c.Scope(scopes.UserSession)
		scopeService := ioc.Get[scopes.UserSessionService](c)
		sessionClient := ioc.Get[clients.SessionClient](c)
		connection := ioc.Get[netconnection.NetConnection](c).Connect(conn, func() {
			scopeService.Clean(scopes.NewUserSessionEndArgs())
		})
		client := clients.NewClient(
			clients.ClientID(conn.RemoteAddr().String()),
			connection,
		)
		sessionClient.UseClient(client)
	}
}

type Pkg struct {
	host     string // CONN_HOST = "localhost"
	port     string // CONN_PORT = "8080"
	connType string // CONN_TYPE = "tcp"
}

func Package(
	// e.g. localhost
	host string,
	// e.g. 8080
	port string,
	// e.g. tcp
	connType string,
) Pkg {
	return Pkg{
		host:     host,
		port:     port,
		connType: connType,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, r runtime.Builder) runtime.Builder {
		r.OnStart(func() {
			go pkg.hostServer(c)
		})
		return r
	})
}

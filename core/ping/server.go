package ping

import (
	"backend/services/clients"
	"fmt"
	"shared/services/logger"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type serverPingEndpoint struct {
	Logger logger.Logger         `inject:"1"`
	Client clients.SessionClient `inject:"1"`
}

func (e serverPingEndpoint) Handle(req PingReq) (PingRes, error) {
	e.Logger.Info(fmt.Sprintf("pinged backend: %v", req))
	client, err := e.Client.Client()
	if err != nil {
		e.Logger.Error(err)
		return PingRes{ID: req.ID, Ok: false}, err
	}
	res, err := relay.Handle(client.Connection.Relay(), PingReq{ID: req.ID + 10})
	e.Logger.Info(fmt.Sprintf("client responsed: \nres is: %v\nerr is: %s\n", res, err))
	return PingRes{ID: req.ID, Ok: true}, nil
}

type ServerPkg struct{}

func ServerPackage() ServerPkg {
	return ServerPkg{}
}

func (pkg ServerPkg) Register(b ioc.Builder) {
	Package().Register(b)
	endpoint.Register[serverPingEndpoint](b)
}

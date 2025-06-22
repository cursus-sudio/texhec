package ping

import (
	"backend/services/api"
	"backend/services/logger"
	"backend/utils/endpoint"
	"fmt"
	"frontend/services/console"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type PingRes struct {
	ID int
	Ok bool
}

type PingReq struct {
	endpoint.Request[PingRes]
	ID int
}

type serverPingEndpoint struct {
	Logger logger.Logger `inject:"1"`
	Client api.Client    `inject:"1"`
}

func (endpoint serverPingEndpoint) Handle(req PingReq) (PingRes, error) {
	endpoint.Logger.Info("aok")
	res, err := relay.Handle(endpoint.Client.Relay(), PingReq{ID: req.ID + 10})
	endpoint.Logger.Info(fmt.Sprintf("server recieved: \nres is: %v\nerr is: %s\n", res, err))
	return PingRes{ID: req.ID, Ok: true}, nil
}

type clientPingEndpoint struct {
	Console console.Console `inject:"1"`
}

func (endpoint clientPingEndpoint) Handle(req PingReq) (PingRes, error) {
	endpoint.Console.LogToConsole(fmt.Sprintf("\nclient endpoint recieved: %v\n", req))
	return PingRes{ID: req.ID, Ok: true}, nil
}

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s api.ServerBuilder) api.ServerBuilder {
		endpoint.Register[serverPingEndpoint](s.Relay())
		return s
	})
}

type ClientPkg struct{}

func ClientPackage() ClientPkg {
	return ClientPkg{}
}

func (pkg ClientPkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s api.ClientBuilder) api.ClientBuilder {
		relay.Register(s.Relay(), func(req PingReq) (PingRes, error) {
			return ioc.GetServices[clientPingEndpoint](c).Handle(req)
		})
		return s
	})
}

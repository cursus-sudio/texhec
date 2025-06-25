package ping

import (
	"fmt"
	"frontend/services/console"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type clientPingEndpoint struct {
	Console console.Console `inject:"1"`
}

func (endpoint clientPingEndpoint) Handle(req PingReq) (PingRes, error) {
	endpoint.Console.LogToConsole(fmt.Sprintf("\nclient endpoint recieved: %v\n", req))
	return PingRes{ID: req.ID, Ok: true}, nil
}

type ClientPkg struct{}

func ClientPackage() ClientPkg {
	return ClientPkg{}
}

func (pkg ClientPkg) Register(b ioc.Builder) {
	Package().Register(b)
	endpoint.Register[clientPingEndpoint](b)
}

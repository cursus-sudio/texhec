package ping

import (
	"fmt"
	"frontend/services/console"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type frontendPingEndpoint struct {
	Console console.Console `inject:"1"`
}

func (endpoint frontendPingEndpoint) Handle(req PingReq) (PingRes, error) {
	endpoint.Console.Print(fmt.Sprintf("\npinged client: %v\n", req))
	return PingRes{ID: req.ID, Ok: true}, nil
}

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (pkg FrontendPkg) Register(b ioc.Builder) {
	Package().Register(b)
	endpoint.Register[frontendPingEndpoint](b)
}

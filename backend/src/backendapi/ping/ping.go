package ping

import (
	"backend/src/backendapi"
	"backend/src/utils/endpoint"
	"backend/src/utils/logger"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type PingRes struct {
	ID int
	Ok bool
}

type PingReq struct {
	relay.Req[PingRes]
	ID int
}

type pingEndpoint struct {
	Logger logger.Logger `inject:"1"`
}

func (endpoint pingEndpoint) Handle(req PingReq) (PingRes, error) {
	endpoint.Logger.Info("aok")
	return PingRes{ID: req.ID, Ok: true}, nil
}

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s backendapi.Builder) backendapi.Builder {
		endpoint.Register[pingEndpoint](c, s.Relay())
		return s
	})
}

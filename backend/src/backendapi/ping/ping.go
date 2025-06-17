package ping

import (
	"backend/src/utils/endpoint"
	"backend/src/utils/logger"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/relay"
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

type Pkg struct {
	c ioc.Dic
}

func Package(c ioc.Dic) Pkg {
	return Pkg{
		c: c,
	}
}

func (pkg Pkg) Register(r relay.Relay) {
	endpoint.Register[pingEndpoint](pkg.c, r)
}

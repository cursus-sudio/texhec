package tacticalmap

import (
	"backend/services/backendapi"
	"backend/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) TacticalMap { return newTacticalMap() })
	ioc.WrapService[backendapi.Builder](b, func(c ioc.Dic, s backendapi.Builder) backendapi.Builder {
		r := s.Relay()
		endpoint.Register[createEndpoint](c, r)
		endpoint.Register[destroyEndpoint](c, r)
		endpoint.Register[getEndpoint](c, r)
		return s
	})
}

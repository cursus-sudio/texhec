package ping

import (
	"reflect"
	"shared/services/requestcodec"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type PingRes struct {
	ID int
	Ok bool
}

type PingReq struct {
	endpoint.Request[PingRes]
	ID int
}

//

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s requestcodec.Builder) requestcodec.Builder {
		s.Register(reflect.TypeFor[PingRes]())
		s.Register(reflect.TypeFor[PingReq]())
		return s
	})
}

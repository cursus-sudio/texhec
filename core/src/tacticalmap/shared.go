package tacticalmap

import (
	"reflect"
	"shared/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s codec.Builder) codec.Builder {
		s.Register(reflect.TypeFor[CreatedMessage]())
		s.Register(reflect.TypeFor[DestroyedMessage]())

		s.Register(reflect.TypeFor[CreateReq]())
		s.Register(reflect.TypeFor[CreateRes]())

		s.Register(reflect.TypeFor[DestroyReq]())
		s.Register(reflect.TypeFor[DestroyRes]())

		s.Register(reflect.TypeFor[GetReq]())
		s.Register(reflect.TypeFor[GetRes]())

		return s
	})
}

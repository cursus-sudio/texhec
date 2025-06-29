package netconnection

import (
	"reflect"
	"shared/services/codec"
	"shared/services/uuid"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		b.Register(reflect.TypeFor[Msg]())
		return b
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) NetConnection {
		return newNetConnection(
			ioc.Get[codec.Codec](c),
			ioc.Get[connection.Connection](c),
			ioc.Get[uuid.Factory](c),
		)
	})
}

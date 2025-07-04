package netconnection

import (
	"reflect"
	"shared/services/codec"
	"shared/services/uuid"
	"shared/utils/connection"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	timeout time.Duration
}

func Package(timeout time.Duration) Pkg {
	return Pkg{
		timeout: timeout,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		b.Register(reflect.TypeFor[Msg]())
		return b
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) NetConnection {
		return newNetConnection(
			ioc.Get[codec.Codec](c),
			ioc.Get[connection.Connection](c),
			ioc.Get[uuid.Factory](c),
			pkg.timeout,
		)
	})

	ioc.RegisterDependency[NetConnection, codec.Codec](b)
	ioc.RegisterDependency[NetConnection, connection.Connection](b)
	ioc.RegisterDependency[NetConnection, uuid.Factory](b)
}

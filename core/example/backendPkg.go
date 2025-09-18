package example

import (
	"backend/services/saves"
	"reflect"
	"shared/services/codec"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type BackendPkg struct{}

func BackendPackage() BackendPkg {
	return BackendPkg{}
}

func (BackendPkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s saves.WorldStateCodecBuilder) saves.WorldStateCodecBuilder {
		saves.AddPersistedArray[IntComponent](s)
		return s
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, w ecs.World) ecs.World {
		e := w.NewEntity()
		ecs.SaveComponent(w.Components(), e, IntComponent{2137})
		return w
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		b.Register(reflect.TypeFor[IntComponent]())
		return b
	})
}

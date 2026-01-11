package test

import (
	"engine/modules/record"
	recordpkg "engine/modules/record/pkg"
	"engine/modules/uuid"
	uuidpkg "engine/modules/uuid/pkg"
	"engine/services/clock"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Component struct {
	Counter int
}

type Setup struct {
	Config record.Config
	Codec  codec.Codec

	World
	ComponentArray ecs.ComponentsArray[Component]
}

type World interface {
	ecs.World
	uuid.UUIDTool
	record.RecordTool
}

type world struct {
	ecs.World
	uuid.UUIDTool
	record.RecordTool
}

func NewSetup() Setup {
	b := ioc.NewBuilder()

	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		codec.Package(),
		uuidpkg.Package(),
		recordpkg.Package(),
	} {
		pkg.Register(b)
	}

	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			Register(Component{})
	})

	c := b.Build()

	w := world{World: ecs.NewWorld()}
	w.UUIDTool = ioc.Get[uuid.ToolFactory](c).Build(w)
	w.RecordTool = ioc.Get[record.ToolFactory](c).Build(w)

	s := Setup{
		Codec:  ioc.Get[codec.Codec](c),
		Config: record.NewConfig(),

		World:          w,
		ComponentArray: ecs.GetComponentsArray[Component](w),
	}

	record.AddToConfig[Component](s.Config)

	return s
}

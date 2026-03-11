package test

import (
	"engine/modules/registry"
	registrypkg "engine/modules/registry/pkg"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	World   ecs.World        `inject:"1"`
	Service registry.Service `inject:"1"`
}

type TagValueComponent struct {
	Value string
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	ecs.Package().Register(b)
	registrypkg.Package().Register(b)
	ioc.WrapService(b, func(c ioc.Dic, registry registry.Service) {
		world := ioc.Get[ecs.World](c)
		registry.Register("tag", func(structTagValue string) ecs.EntityID {
			entity := world.NewEntity()
			ecs.SaveComponent(world, entity, TagValueComponent{structTagValue})
			return entity
		})
	})
	c := b.Build()
	return ioc.GetServices[Setup](c)
}

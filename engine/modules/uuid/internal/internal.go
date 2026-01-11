package internal

import (
	"engine/modules/relation"
	"engine/modules/uuid"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	uuidArray ecs.ComponentsArray[uuid.Component]
	relation.Service[uuid.UUID]
	uuid.Factory
}

func NewService(c ioc.Dic) uuid.Service {
	t := &service{
		ecs.GetComponentsArray[uuid.Component](ioc.Get[ecs.World](c)),
		ioc.Get[relation.Service[uuid.UUID]](c),
		ioc.Get[uuid.Factory](c),
	}
	return t
}

func (t *service) UUID() ecs.ComponentsArray[uuid.Component] { return t.uuidArray }

func (t *service) Entity(uuid uuid.UUID) (ecs.EntityID, bool) {
	return t.Get(uuid)
}

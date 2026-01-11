package internal

import (
	"engine/modules/relation"
	"engine/modules/uuid"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World                       ecs.World `inject:"1"`
	relation.Service[uuid.UUID] `inject:"1"`
	uuid.Factory                `inject:"1"`

	uuidArray ecs.ComponentsArray[uuid.Component]
}

func NewService(c ioc.Dic) uuid.Service {
	t := ioc.GetServices[*service](c)
	t.uuidArray = ecs.GetComponentsArray[uuid.Component](t.World)
	return t
}

func (t *service) Component() ecs.ComponentsArray[uuid.Component] { return t.uuidArray }

func (t *service) Entity(uuid uuid.UUID) (ecs.EntityID, bool) {
	return t.Get(uuid)
}

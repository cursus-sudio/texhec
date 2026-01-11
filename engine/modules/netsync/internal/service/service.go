package service

import (
	"engine/modules/netsync"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World  ecs.World `inject:"1"`
	server ecs.ComponentsArray[netsync.ServerComponent]
	client ecs.ComponentsArray[netsync.ClientComponent]
}

func NewService(c ioc.Dic) netsync.Service {
	t := ioc.GetServices[*service](c)
	t.server = ecs.GetComponentsArray[netsync.ServerComponent](t.World)
	t.client = ecs.GetComponentsArray[netsync.ClientComponent](t.World)
	return t
}

func (t *service) Server() ecs.ComponentsArray[netsync.ServerComponent] { return t.server }
func (t *service) Client() ecs.ComponentsArray[netsync.ClientComponent] { return t.client }

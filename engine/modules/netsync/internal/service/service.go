package service

import (
	"engine/modules/netsync"
	"engine/services/ecs"
)

type service struct {
	server ecs.ComponentsArray[netsync.ServerComponent]
	client ecs.ComponentsArray[netsync.ClientComponent]
}

func NewService(
	world ecs.World,
) netsync.Service {
	t := &service{
		server: ecs.GetComponentsArray[netsync.ServerComponent](world),
		client: ecs.GetComponentsArray[netsync.ClientComponent](world),
	}
	return t
}

func (t *service) Server() ecs.ComponentsArray[netsync.ServerComponent] { return t.server }
func (t *service) Client() ecs.ComponentsArray[netsync.ClientComponent] { return t.client }

package service

import (
	"engine/modules/netsync"
	"engine/services/ecs"
)

type tool struct {
	server ecs.ComponentsArray[netsync.ServerComponent]
	client ecs.ComponentsArray[netsync.ClientComponent]
}

func NewToolFactory(
	world ecs.World,
) netsync.Service {
	t := &tool{
		server: ecs.GetComponentsArray[netsync.ServerComponent](world),
		client: ecs.GetComponentsArray[netsync.ClientComponent](world),
	}
	return t
}

func (t *tool) Server() ecs.ComponentsArray[netsync.ServerComponent] { return t.server }
func (t *tool) Client() ecs.ComponentsArray[netsync.ClientComponent] { return t.client }

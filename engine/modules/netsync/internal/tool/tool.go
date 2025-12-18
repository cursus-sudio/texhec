package tool

import (
	"engine/modules/netsync"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	server ecs.ComponentsArray[netsync.ServerComponent]
	client ecs.ComponentsArray[netsync.ClientComponent]
}

func NewToolFactory() ecs.ToolFactory[netsync.World, netsync.NetSyncTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w netsync.World) netsync.NetSyncTool {
		mutex.Lock()
		defer mutex.Unlock()

		t := tool{
			server: ecs.GetComponentsArray[netsync.ServerComponent](w),
			client: ecs.GetComponentsArray[netsync.ClientComponent](w),
		}
		return t
	})
}

func (t tool) NetSync() netsync.Interface {
	return t
}

func (t tool) Server() ecs.ComponentsArray[netsync.ServerComponent] { return t.server }
func (t tool) Client() ecs.ComponentsArray[netsync.ClientComponent] { return t.client }

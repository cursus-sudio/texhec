package server

import (
	"engine/modules/netsync"
	"engine/modules/netsync/internal/config"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

func NewToolFactory(
	config config.Config,
	netSyncToolFactory netsync.ToolFactory,
	logger logger.Logger,
) ecs.ToolFactory[netsync.World, *Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w netsync.World) *Tool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[Tool](w); ok {
			return t
		}
		t := NewTool(
			config,
			netSyncToolFactory,
			logger,
			w,
		)
		w.SaveGlobal(t)
		return t
	})
}

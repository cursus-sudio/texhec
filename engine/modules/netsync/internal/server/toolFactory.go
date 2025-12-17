package server

import (
	"engine/modules/netsync/internal/config"
	"engine/modules/netsync/internal/state"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

func NewToolFactory(
	config config.Config,
	stateToolFactory ecs.ToolFactory[state.Tool],
	uniqueToolFactory ecs.ToolFactory[uuid.UUIDTool],
	logger logger.Logger,
) ecs.ToolFactory[Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) Tool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[Tool](w); ok {
			return t
		}
		t := NewTool(
			config,
			uniqueToolFactory,
			stateToolFactory,
			logger,
			w,
		)
		w.SaveGlobal(t)
		return t
	})
}

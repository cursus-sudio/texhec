package server

import (
	"engine/modules/sync/internal/config"
	"engine/modules/sync/internal/state"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

func NewToolFactory(
	config config.Config,
	stateToolFactory ecs.ToolFactory[state.Tool],
	uniqueToolFactory ecs.ToolFactory[uuid.Tool],
	logger logger.Logger,
) ecs.ToolFactory[Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) Tool {
		if t, err := ecs.GetGlobal[Tool](w); err == nil {
			return t
		}
		mutex.Lock()
		defer mutex.Unlock()
		if t, err := ecs.GetGlobal[Tool](w); err == nil {
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

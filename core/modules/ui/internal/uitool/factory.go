package uitool

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
	"time"
)

func NewToolFactory(
	animationDuration time.Duration,
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
) ui.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ui.World) ui.UiTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := NewTool(
			animationDuration,
			w,
			gameAssets,
			logger,
		)
		w.SaveGlobal(t)
		return t
	})
}

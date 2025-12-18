package uitool

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
	"time"
)

func NewToolFactory(
	animationDuration time.Duration,
	showAnimation animation.AnimationID,
	hideAnimation animation.AnimationID,
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
) ecs.ToolFactory[ui.World, ui.UiTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ui.World) ui.UiTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := NewTool(
			animationDuration,
			showAnimation,
			hideAnimation,
			w,
			gameAssets,
			logger,
		)
		w.SaveGlobal(t)
		return t
	})
}

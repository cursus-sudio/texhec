package uitool

import (
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/hierarchy"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
	"time"
)

func NewToolFactory(
	animationDuration time.Duration,
	showAnimation animation.AnimationID,
	hideAnimation animation.AnimationID,
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
	tileToolFactory ecs.ToolFactory[tile.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
	renderToolFactory ecs.ToolFactory[render.Tool],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.Tool],
) ecs.ToolFactory[ui.Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) ui.Tool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, err := ecs.GetGlobal[tool](w); err == nil {
			return t
		}
		t := NewTool(
			animationDuration,
			showAnimation,
			hideAnimation,
			w,
			logger,
			cameraToolFactory,
			transformToolFactory,
			tileToolFactory,
			textToolFactory,
			renderToolFactory,
			hierarchyToolFactory,
		)
		w.SaveGlobal(t)
		return t
	})
}

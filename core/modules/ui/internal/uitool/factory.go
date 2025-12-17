package uitool

import (
	gameassets "core/assets"
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
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.CameraTool],
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	tileToolFactory ecs.ToolFactory[tile.Tile],
	textToolFactory ecs.ToolFactory[text.TextTool],
	renderToolFactory ecs.ToolFactory[render.RenderTool],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.HierarchyTool],
) ecs.ToolFactory[ui.UiTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) ui.UiTool {
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

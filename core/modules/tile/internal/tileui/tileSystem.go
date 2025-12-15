package tileui

import (
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/text"
	"engine/services/ecs"
	"engine/services/logger"
	"fmt"

	"github.com/ogiusek/events"
)

func NewSystem(
	logger logger.Logger,
	uiToolFactory ecs.ToolFactory[ui.Tool],
	textToolFactory ecs.ToolFactory[text.Text],
	tileToolFactory ecs.ToolFactory[tile.Tile],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](world)
		uiTool := uiToolFactory.Build(world)
		textTool := textToolFactory.Build(world)
		tileTool := tileToolFactory.Build(world)

		events.Listen(world.EventsBuilder(), func(e tile.TileClickEvent) {
			entity, ok := tileTool.Tile().TilePos().Get(e.Tile)
			if !ok {
				logger.Warn(fmt.Errorf("entity with uuid should exist"))
				return
			}
			pos, ok := tilePosArray.GetComponent(entity)
			if !ok {
				return
			}
			p := uiTool.Show()
			textTool.Text().TextContent().SaveComponent(p, text.TextComponent{Text: fmt.Sprintf("TILE: %v", pos)})
			textTool.Text().FontSize().SaveComponent(p, text.FontSizeComponent{FontSize: 25})
			textTool.Text().TextAlign().SaveComponent(p, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		})
		return nil
	})
}

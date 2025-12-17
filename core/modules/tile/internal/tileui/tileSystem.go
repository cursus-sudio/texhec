package tileui

import (
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"fmt"

	"github.com/ogiusek/events"
)

func NewSystem(
	logger logger.Logger,
	uiToolFactory ecs.ToolFactory[ui.UiTool],
	textToolFactory ecs.ToolFactory[text.TextTool],
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.HierarchyTool],
	tileToolFactory ecs.ToolFactory[tile.Tile],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](world)
		uiTool := uiToolFactory.Build(world).Ui()
		textTool := textToolFactory.Build(world).Text()
		transformTool := transformToolFactory.Build(world).Transform()
		hierarchyTool := hierarchyToolFactory.Build(world).Hierarchy()
		tileTool := tileToolFactory.Build(world).Tile()
		inheritGroupsArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](world)

		events.Listen(world.EventsBuilder(), func(e tile.TileClickEvent) {
			tileEntity, ok := tileTool.TilePos().Get(e.Tile)
			if !ok {
				logger.Warn(fmt.Errorf("entity with uuid should exist"))
				return
			}
			pos, ok := tilePosArray.Get(tileEntity)
			if !ok {
				return
			}
			p := uiTool.Show()
			entity := world.NewEntity()
			hierarchyTool.SetParent(entity, p)
			transformTool.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXYZ))
			inheritGroupsArray.Set(entity, groups.InheritGroupsComponent{})

			textTool.Content().Set(entity, text.TextComponent{Text: fmt.Sprintf("TILE: %v", pos)})
			textTool.FontSize().Set(entity, text.FontSizeComponent{FontSize: 25})
			textTool.Align().Set(entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		})
		return nil
	})
}

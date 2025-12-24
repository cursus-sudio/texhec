package tileui

import (
	"core/modules/tile"
	"engine/modules/groups"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"fmt"

	"github.com/ogiusek/events"
)

func NewSystem(
	logger logger.Logger,
	tileToolFactory tile.ToolFactory,
) tile.System {
	return ecs.NewSystemRegister(func(world tile.World) error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](world)
		tileTool := tileToolFactory.Build(world).Tile()
		inheritGroupsArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](world)

		events.Listen(world.EventsBuilder(), func(e tile.TileClickEvent) {
			tileEntity, ok := tileTool.PosKey().Get(e.Tile)
			if !ok {
				logger.Warn(fmt.Errorf("entity with uuid should exist"))
				return
			}
			pos, ok := tilePosArray.Get(tileEntity)
			if !ok {
				return
			}
			p := world.Ui().Show()
			entity := world.NewEntity()
			world.Hierarchy().SetParent(entity, p)
			world.Transform().Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXYZ))
			inheritGroupsArray.Set(entity, groups.InheritGroupsComponent{})

			world.Text().Content().Set(entity, text.TextComponent{Text: fmt.Sprintf("TILE: %v", pos)})
			world.Text().FontSize().Set(entity, text.FontSizeComponent{FontSize: 25})
			world.Text().Align().Set(entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		})
		return nil
	})
}

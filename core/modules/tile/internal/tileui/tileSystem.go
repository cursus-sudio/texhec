package tileui

import (
	"core/modules/tile"
	"core/modules/ui"
	"engine"
	"engine/modules/groups"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"fmt"

	"github.com/ogiusek/events"
)

func NewSystem(
	world engine.World,
	ui ui.Service,
	tileService tile.Service,
) tile.System {
	return ecs.NewSystemRegister(func() error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](world)
		inheritGroupsArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](world)

		events.Listen(world.EventsBuilder, func(e tile.TileClickEvent) {
			tileEntity, ok := tileService.PosKey().Get(e.Tile)
			if !ok {
				world.Logger.Warn(fmt.Errorf("entity with uuid should exist"))
				return
			}
			pos, ok := tilePosArray.Get(tileEntity)
			if !ok {
				return
			}
			p := ui.Show()
			entity := world.NewEntity()
			world.Hierarchy.SetParent(entity, p)
			world.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXYZ))
			inheritGroupsArray.Set(entity, groups.InheritGroupsComponent{})

			world.Text.Content().Set(entity, text.TextComponent{Text: fmt.Sprintf("TILE: %v", pos)})
			world.Text.FontSize().Set(entity, text.FontSizeComponent{FontSize: 25})
			world.Text.Align().Set(entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		})
		return nil
	})
}

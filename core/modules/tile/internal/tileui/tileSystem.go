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
	"github.com/ogiusek/ioc/v2"
)

type system struct {
	engine.World `inject:"1"`
	Ui           ui.Service   `inject:"1"`
	Tile         tile.Service `inject:"1"`
}

func NewSystem(c ioc.Dic) tile.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[system](c)
		inheritGroupsArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](s)

		events.Listen(s.EventsBuilder, func(e tile.TileClickEvent) {
			grid, ok := s.Tile.Grid().Get(e.Grid)
			if !ok {
				s.Logger.Warn(fmt.Errorf("grid doesn't exist"))
				return
			}
			coords := grid.GetCoords(e.Tile)
			if !ok {
				return
			}
			p := s.Ui.Show()
			entity := s.NewEntity()
			s.Hierarchy.SetParent(entity, p)
			s.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXYZ))
			inheritGroupsArray.Set(entity, groups.InheritGroupsComponent{})

			s.Text.Content().Set(entity, text.TextComponent{Text: fmt.Sprintf("TILE: %v", coords)})
			s.Text.FontSize().Set(entity, text.FontSizeComponent{FontSize: 25})
			s.Text.Align().Set(entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		})
		return nil
	})
}

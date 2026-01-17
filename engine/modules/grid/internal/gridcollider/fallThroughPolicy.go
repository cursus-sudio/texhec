package gridcollider

import (
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type ClickEvent[Tile grid.TileConstraint] struct {
	Target inputs.Target
}

func (ClickEvent[Tile]) SetTarget(target inputs.Target) inputs.EventTargetSetter {
	return ClickEvent[Tile]{target}
}

type squareFallThroughPolicy[Tile grid.TileConstraint] struct {
	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`
	World         ecs.World      `inject:"1"`
	Inputs        inputs.Service `inject:"1"`
	Logger        logger.Logger  `inject:"1"`

	DirtyEntities ecs.DirtySet
	GridArray     ecs.ComponentsArray[grid.SquareGridComponent[Tile]]

	indexEvent func(grid.Index) any
}

func NewColliderWithPolicy[Tile grid.TileConstraint](
	c ioc.Dic,
	indexEvent func(grid.Index) any,
) collider.FallTroughPolicy {
	s := ioc.GetServices[*squareFallThroughPolicy[Tile]](c)

	s.GridArray = ecs.GetComponentsArray[grid.SquareGridComponent[Tile]](s.World)
	s.indexEvent = indexEvent

	s.GridArray.AddDirtySet(s.DirtyEntities)
	s.GridArray.BeforeGet(s.BeforeGet)

	events.Listen(s.EventsBuilder, s.OnClick)

	return s
}

func (t *squareFallThroughPolicy[Tile]) BeforeGet() {
	for _, entity := range t.DirtyEntities.Get() {
		if !t.World.EntityExists(entity) {
			continue
		}
		t.Inputs.LeftClick().Set(entity, inputs.NewLeftClick(ClickEvent[Tile]{}))
	}
}

func (t *squareFallThroughPolicy[Tile]) getIndex(
	gridComponent grid.SquareGridComponent[Tile],
	collision collider.ObjectRayCollision,
) (grid.Index, bool) {
	w := float32(gridComponent.Width())
	h := float32(gridComponent.Height())

	point := collision.Hit.Point
	x := grid.Coord(w * (1 + point.X()) / 2)
	y := grid.Coord(h * (1 + point.Y()) / 2)

	index, ok := gridComponent.GetIndex(x, y)
	if !ok {
		return 0, false
	}
	return index, true
}

func (t *squareFallThroughPolicy[Tile]) FallThrough(collision collider.ObjectRayCollision) bool {
	gridComponent, ok := t.GridArray.Get(collision.Entity)
	if !ok {
		return false
	}

	index, ok := t.getIndex(gridComponent, collision)
	if !ok {
		return true
	}

	tile := gridComponent.GetTile(index)
	return tile == 0
}

func (t *squareFallThroughPolicy[Tile]) OnClick(e ClickEvent[Tile]) {
	gridComponent, ok := t.GridArray.Get(e.Target.Entity)
	if !ok {
		return
	}
	index, ok := t.getIndex(gridComponent, e.Target.ObjectRayCollision)
	if !ok {
		return
	}
	event := t.indexEvent(index)
	events.Emit(t.Events, event)

}

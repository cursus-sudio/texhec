package internal

import (
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"golang.org/x/exp/constraints"
)

type ClickEvent[TileType constraints.Unsigned] struct {
	Target inputs.Target
}

func (ClickEvent[TileType]) SetTarget(target inputs.Target) inputs.EventTargetSetter {
	return ClickEvent[TileType]{target}
}

type squareFallThroughPolicy[TileType constraints.Unsigned] struct {
	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`
	World         ecs.World      `inject:"1"`
	Inputs        inputs.Service `inject:"1"`
	Logger        logger.Logger  `inject:"1"`

	DirtyEntities ecs.DirtySet
	GridArray     ecs.ComponentsArray[grid.SquareGridComponent[TileType]]

	indexEvent func(grid.Index) any
}

func NewColliderWithPolicy[TileType constraints.Unsigned](
	c ioc.Dic,
	indexEvent func(grid.Index) any,
) collider.FallTroughPolicy {
	s := ioc.GetServices[*squareFallThroughPolicy[TileType]](c)

	s.GridArray = ecs.GetComponentsArray[grid.SquareGridComponent[TileType]](s.World)
	s.indexEvent = indexEvent

	s.GridArray.AddDirtySet(s.DirtyEntities)
	s.GridArray.BeforeGet(s.BeforeGet)

	events.Listen(s.EventsBuilder, s.OnClick)

	return s
}

func (t *squareFallThroughPolicy[TileType]) BeforeGet() {
	for _, entity := range t.DirtyEntities.Get() {
		if !t.World.EntityExists(entity) {
			continue
		}
		t.Inputs.LeftClick().Set(entity, inputs.NewLeftClick(ClickEvent[TileType]{}))
	}
}

func (t *squareFallThroughPolicy[TileType]) getIndex(
	gridComponent grid.SquareGridComponent[TileType],
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

func (t *squareFallThroughPolicy[TileType]) FallThrough(collision collider.ObjectRayCollision) bool {
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

func (t *squareFallThroughPolicy[TileType]) OnClick(e ClickEvent[TileType]) {
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

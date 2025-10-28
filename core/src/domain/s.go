package domain

import (
	gameassets "core/assets"
	"core/src/tile"
	"fmt"
	"frontend/engine/components/collider"
	"frontend/engine/components/mouse"
	"frontend/services/console"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/ogiusek/events"
)

type QuitEvent struct{}

func NewQuitEvent() QuitEvent { return QuitEvent{} }

//

type OnHoveredEvent struct {
	entity   ecs.EntityID
	row, col int
}

func NewOnHoveredEvent(entity ecs.EntityID, row, col int) OnHoveredEvent {
	return OnHoveredEvent{
		entity: entity,
		row:    row, col: col,
	}
}

//

type OnClickEvent struct {
	entity   ecs.EntityID
	row, col int
}

func NewOnClickEvent(entity ecs.EntityID, row, col int) OnClickEvent {
	return OnClickEvent{
		entity: entity,
		row:    row, col: col,
	}
}

//

type sys struct {
	world   ecs.World
	logger  logger.Logger
	runtime runtime.Runtime
	console console.Console

	tileArray        ecs.ComponentsArray[tile.TileComponent]
	colliderArray    ecs.ComponentsArray[collider.Collider]
	mouseEventsArray ecs.ComponentsArray[mouse.MouseEvents]

	colliderTransaction    ecs.ComponentsArrayTransaction[collider.Collider]
	mouseEventsTransaction ecs.ComponentsArrayTransaction[mouse.MouseEvents]
}

func NewSys(
	world ecs.World,
	logger logger.Logger,
	runtime runtime.Runtime,
	console console.Console,
) ecs.SystemRegister {
	tileArray := ecs.GetComponentsArray[tile.TileComponent](world.Components())
	colliderArray := ecs.GetComponentsArray[collider.Collider](world.Components())
	mouseEventsArray := ecs.GetComponentsArray[mouse.MouseEvents](world.Components())
	return &sys{
		world:   world,
		logger:  logger,
		runtime: runtime,
		console: console,

		tileArray:        tileArray,
		colliderArray:    colliderArray,
		mouseEventsArray: mouseEventsArray,

		colliderTransaction:    colliderArray.Transaction(),
		mouseEventsTransaction: mouseEventsArray.Transaction(),
	}
}

func (s *sys) Register(b events.Builder) {
	ecs.RegisterSystems(b,
		ecs.NewSystemRegister(func(b events.Builder) {
			onChangeOrAdd := func(ei []ecs.EntityID) {
				colliderTransaction := s.colliderArray.Transaction()
				mouseEventsTransaction := s.mouseEventsArray.Transaction()
				for _, entity := range ei {
					tile, err := s.tileArray.GetComponent(entity)
					if err != nil {
						continue
					}

					colliderTransaction.SaveComponent(entity, collider.NewCollider(gameassets.SquareColliderID))
					mouseEventsTransaction.SaveComponent(entity, mouse.NewMouseEvents().
						AddLeftClickEvents(OnClickEvent{entity, int(tile.Pos.X), int(tile.Pos.Y)}).
						AddMouseHoverEvents(OnHoveredEvent{entity, int(tile.Pos.X), int(tile.Pos.Y)}),
					)
				}
				err := ecs.FlushMany(colliderTransaction, mouseEventsTransaction)
				if err != nil {
					s.logger.Error(err)
				}
			}

			s.tileArray.OnAdd(onChangeOrAdd)
			s.tileArray.OnChange(onChangeOrAdd)
		}),
		ecs.NewSystemRegister(func(b events.Builder) {
			events.Listen(b, func(e QuitEvent) {
				s.runtime.Stop()
			})
			events.Listen(b, func(e OnHoveredEvent) {
				s.console.Print(
					fmt.Sprintf("damn it really is hovered %v (%d, %d)\n", e.entity, e.col, e.row),
				)
			})
			events.Listen(b, func(e OnClickEvent) {
				s.console.PrintPermanent(
					fmt.Sprintf("damn it really is clicked %v (%d, %d)\n", e.entity, e.col, e.row),
				)
			})
		}),
	)
}

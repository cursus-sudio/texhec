package projections

import (
	"frontend/engine/components/projection"
	"frontend/services/ecs"
	"frontend/services/media/window"
)

// events

// updates dynamic projections
type UpdateProjectionsEvent struct{}

func NewUpdateProjectionsEvent() UpdateProjectionsEvent {
	return UpdateProjectionsEvent{}
}

// system

type UpdateProjetionsSystem struct {
	world  ecs.World
	window window.Api

	perspectivesQuery ecs.LiveQuery
	orthoQuery        ecs.LiveQuery
}

func NewUpdateProjectionsSystem(world ecs.World, window window.Api) UpdateProjetionsSystem {
	perspectiveQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicPerspective{}))
	orthoQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicOrtho{}))
	return UpdateProjetionsSystem{
		world:  world,
		window: window,

		perspectivesQuery: perspectiveQuery,
		orthoQuery:        orthoQuery,
	}
}

func (s UpdateProjetionsSystem) Listen(e UpdateProjectionsEvent) {
	var w, h float32
	{
		width, height := s.window.Window().GetSize()
		w, h = float32(width), float32(height)
	}
	aspectRatio := w / h
	for _, entity := range s.perspectivesQuery.Entities() {
		resizePerspective, err := ecs.GetComponent[projection.DynamicPerspective](s.world, entity)
		if err != nil {
			continue
		}
		perspective := projection.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		s.world.SaveComponent(entity, perspective)
	}
	for _, entity := range s.orthoQuery.Entities() {
		resizeOrtho, err := ecs.GetComponent[projection.DynamicOrtho](s.world, entity)
		if err != nil {
			continue
		}
		ortho := projection.NewOrtho(
			w, h,
			resizeOrtho.Near, resizeOrtho.Far,
		)
		s.world.SaveComponent(entity, ortho)
	}
}

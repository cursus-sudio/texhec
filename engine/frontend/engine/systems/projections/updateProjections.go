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
}

func NewUpdateProjectionsSystem(world ecs.World, window window.Api) UpdateProjetionsSystem {
	return UpdateProjetionsSystem{
		world:  world,
		window: window,
	}
}

func (s UpdateProjetionsSystem) Listen(e UpdateProjectionsEvent) {
	var w, h float32
	{
		width, height := s.window.Window().GetSize()
		w, h = float32(width), float32(height)
	}
	aspectRatio := w / h

	for _, entity := range s.world.GetEntitiesWithComponents(ecs.GetComponentType(projection.DynamicPerspective{})) {
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
	for _, entity := range s.world.GetEntitiesWithComponents(ecs.GetComponentType(projection.DynamicOrtho{})) {
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

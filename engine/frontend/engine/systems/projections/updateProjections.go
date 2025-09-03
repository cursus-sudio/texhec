package projections

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/ecs"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
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
	s := UpdateProjetionsSystem{
		world:  world,
		window: window,

		perspectivesQuery: perspectiveQuery,
		orthoQuery:        orthoQuery,
	}
	listener := func(_ []ecs.EntityID) { s.Listen(UpdateProjectionsEvent{}) }
	perspectiveQuery.OnAdd(listener)
	perspectiveQuery.OnChange(listener)
	perspectiveQuery.OnRemove(listener)
	orthoQuery.OnAdd(listener)
	orthoQuery.OnChange(listener)
	orthoQuery.OnRemove(listener)

	return s
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
		originalPerspective, _ := ecs.GetComponent[projection.Perspective](s.world, entity)
		if originalPerspective == perspective {
			continue
		}
		s.world.SaveComponent(entity, perspective)
	}
	for _, entity := range s.orthoQuery.Entities() {
		transformComponent, err := ecs.GetComponent[transform.Transform](s.world, entity)
		if err != nil {
			continue
		}
		resizeOrtho, err := ecs.GetComponent[projection.DynamicOrtho](s.world, entity)
		if err != nil {
			continue
		}

		transformComponent.Size = mgl32.Vec3{w / resizeOrtho.Zoom, h / resizeOrtho.Zoom, transformComponent.Size.Z()}
		s.world.SaveComponent(entity, transformComponent)

		ortho := projection.NewOrtho(
			w/resizeOrtho.Zoom, h/resizeOrtho.Zoom,
			resizeOrtho.Near, resizeOrtho.Far,
		)
		originalOrtho, _ := ecs.GetComponent[projection.Ortho](s.world, entity)
		if originalOrtho == ortho {
			continue
		}
		s.world.SaveComponent(entity, ortho)
	}
}

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

	transformArray ecs.ComponentsArray[transform.Transform]

	dynamicPerspectivesArray ecs.ComponentsArray[projection.DynamicPerspective]
	dynamicOrthoArray        ecs.ComponentsArray[projection.DynamicOrtho]
	perspectivesArray        ecs.ComponentsArray[projection.Perspective]
	orthoArray               ecs.ComponentsArray[projection.Ortho]
}

func NewUpdateProjectionsSystem(world ecs.World, window window.Api) UpdateProjetionsSystem {
	perspectiveQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicPerspective{}))
	orthoQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicOrtho{}))
	s := UpdateProjetionsSystem{
		world:  world,
		window: window,

		perspectivesQuery: perspectiveQuery,
		orthoQuery:        orthoQuery,

		transformArray: ecs.GetComponentArray[transform.Transform](world.Components()),

		dynamicPerspectivesArray: ecs.GetComponentArray[projection.DynamicPerspective](world.Components()),
		dynamicOrthoArray:        ecs.GetComponentArray[projection.DynamicOrtho](world.Components()),
		perspectivesArray:        ecs.GetComponentArray[projection.Perspective](world.Components()),
		orthoArray:               ecs.GetComponentArray[projection.Ortho](world.Components()),
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
		resizePerspective, err := s.dynamicPerspectivesArray.GetComponent(entity)
		if err != nil {
			continue
		}
		perspective := projection.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		originalPerspective, _ := s.perspectivesArray.GetComponent(entity)
		if originalPerspective == perspective {
			continue
		}
		s.perspectivesArray.SaveComponent(entity, perspective)
	}
	for _, entity := range s.orthoQuery.Entities() {
		transformComponent, err := s.transformArray.GetComponent(entity)
		if err != nil {
			continue
		}
		resizeOrtho, err := s.dynamicOrthoArray.GetComponent(entity)
		if err != nil {
			continue
		}

		transformComponent.Size = mgl32.Vec3{w / resizeOrtho.Zoom, h / resizeOrtho.Zoom, transformComponent.Size.Z()}
		s.transformArray.SaveComponent(entity, transformComponent)

		ortho := projection.NewOrtho(
			w/resizeOrtho.Zoom, h/resizeOrtho.Zoom,
			resizeOrtho.Near, resizeOrtho.Far,
		)
		originalOrtho, _ := s.orthoArray.GetComponent(entity)
		if originalOrtho == ortho {
			continue
		}
		s.orthoArray.SaveComponent(entity, ortho)
	}
}

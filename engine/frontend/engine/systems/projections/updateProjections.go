package projections

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

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
	logger logger.Logger

	perspectivesQuery ecs.LiveQuery
	orthoQuery        ecs.LiveQuery

	transformArray ecs.ComponentsArray[transform.Transform]

	dynamicPerspectivesArray ecs.ComponentsArray[projection.DynamicPerspective]
	dynamicOrthoArray        ecs.ComponentsArray[projection.DynamicOrtho]
	perspectivesArray        ecs.ComponentsArray[projection.Perspective]
	orthoArray               ecs.ComponentsArray[projection.Ortho]
}

func NewUpdateProjectionsSystem(world ecs.World, window window.Api, logger logger.Logger) UpdateProjetionsSystem {
	perspectiveQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicPerspective{}))
	orthoQuery := world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.DynamicOrtho{}))
	s := UpdateProjetionsSystem{
		world:  world,
		window: window,
		logger: logger,

		perspectivesQuery: perspectiveQuery,
		orthoQuery:        orthoQuery,

		transformArray: ecs.GetComponentsArray[transform.Transform](world.Components()),

		dynamicPerspectivesArray: ecs.GetComponentsArray[projection.DynamicPerspective](world.Components()),
		dynamicOrthoArray:        ecs.GetComponentsArray[projection.DynamicOrtho](world.Components()),
		perspectivesArray:        ecs.GetComponentsArray[projection.Perspective](world.Components()),
		orthoArray:               ecs.GetComponentsArray[projection.Ortho](world.Components()),
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
	transformTransaction := s.transformArray.Transaction()
	perspectiveTransaction := s.perspectivesArray.Transaction()
	orthoTransaction := s.orthoArray.Transaction()

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
		perspectiveTransaction.SaveComponent(entity, perspective)
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

		transformComponent.SetSize(mgl32.Vec3{
			w / resizeOrtho.Zoom, h / resizeOrtho.Zoom, transformComponent.Size.Z(),
		})
		transformTransaction.SaveComponent(entity, transformComponent)

		ortho := projection.NewOrtho(
			w/resizeOrtho.Zoom, h/resizeOrtho.Zoom,
			resizeOrtho.Near, resizeOrtho.Far,
		)
		originalOrtho, _ := s.orthoArray.GetComponent(entity)
		if originalOrtho == ortho {
			continue
		}
		orthoTransaction.SaveComponent(entity, ortho)
	}

	if err := ecs.FlushMany(
		transformTransaction,
		perspectiveTransaction,
		orthoTransaction,
	); err != nil {
		s.logger.Error(err)
	}

}

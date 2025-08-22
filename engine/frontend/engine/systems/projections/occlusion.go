package projections

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/datastructures"
	"frontend/services/ecs"
)

type OcclusionSystem struct {
	orthoEntities              datastructures.Set[ecs.EntityID]
	perspectiveEntities        datastructures.Set[ecs.EntityID]
	visibleOrthoEntities       datastructures.Set[ecs.EntityID]
	visiblePerspectiveEntities datastructures.Set[ecs.EntityID]
}

func NewOcclusionSystem(world ecs.World) *OcclusionSystem {
	orthoEntities := datastructures.NewSet[ecs.EntityID]()
	perspectiveEntities := datastructures.NewSet[ecs.EntityID]()

	// visible ortho
	visibleOrthoEntities := datastructures.NewSet[ecs.EntityID]()
	visiblePerspectiveEntities := datastructures.NewSet[ecs.EntityID]()

	s := &OcclusionSystem{
		orthoEntities:              orthoEntities,
		perspectiveEntities:        perspectiveEntities,
		visibleOrthoEntities:       visibleOrthoEntities,
		visiblePerspectiveEntities: visiblePerspectiveEntities,
	}

	// add and remove entities
	entitiesQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(projection.UsedProjection{}),
		ecs.GetComponentType(transform.Transform{}),
	)

	entitiesQuery.OnAdd(func(ei []ecs.EntityID) {
		for _, entity := range ei {
			world.SaveComponent(entity, projection.Visible{})
		}
	})

	// orthoQuery := world.QueryEntitiesWithComponents(
	// 	ecs.GetComponentType(projection.Ortho{}),
	// 	ecs.GetComponentType(transform.Transform{}),
	// )
	//
	// perspectiveQuery := world.QueryEntitiesWithComponents(
	// 	ecs.GetComponentType(projection.Ortho{}),
	// 	ecs.GetComponentType(transform.Transform{}),
	// )

	return s
}

// func (s *OcclusionSystem) occlude() {
// }

package projections

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type orthoOcclusion struct {
	tileSize          mgl32.Vec2
	positionsEntities map[mgl32.Vec2]datastructures.Set[ecs.EntityID]
	entitiesMeta      map[ecs.EntityID]struct{ aabb transform.AABB }
}

func newOrthoOcclusion(tileSize mgl32.Vec2) *orthoOcclusion {
	return &orthoOcclusion{
		tileSize:          tileSize,
		positionsEntities: map[mgl32.Vec2]datastructures.Set[ecs.EntityID]{},
		entitiesMeta:      map[ecs.EntityID]struct{ aabb transform.AABB }{},
	}
}

func floorF32ToInt(num float32) int {
	return int(math.Floor(float64(num)))
}

func (s *orthoOcclusion) getPositions(aabb transform.AABB) []mgl32.Vec2 {
	minPos, maxPos := aabb.Min.Vec2(), aabb.Max.Vec2()
	minGridX := floorF32ToInt(minPos.X() / s.tileSize.X())
	minGridY := floorF32ToInt(minPos.Y() / s.tileSize.Y())
	maxGridX := floorF32ToInt(maxPos.X() / s.tileSize.X())
	maxGridY := floorF32ToInt(maxPos.Y() / s.tileSize.Y())

	var positions []mgl32.Vec2
	for x := minGridX; x <= maxGridX; x++ {
		for y := minGridY; y <= maxGridY; y++ {
			tileX := float32(x) * s.tileSize.X()
			tileY := float32(y) * s.tileSize.Y()
			tileCenter := mgl32.Vec2{tileX, tileY}

			positions = append(positions, tileCenter)
		}
	}

	return positions
}

func (s *orthoOcclusion) Add(entity ecs.EntityID, aabb transform.AABB) {
	positions := s.getPositions(aabb)
	if len(positions) == 0 {
		return
	}
	s.entitiesMeta[entity] = struct{ aabb transform.AABB }{aabb}

	for _, position := range positions {
		set, ok := s.positionsEntities[position]
		if !ok {
			set = datastructures.NewSet[ecs.EntityID]()
			s.positionsEntities[position] = set
		}
		set.Add(entity)
	}
}

func (s *orthoOcclusion) Remove(entity ecs.EntityID) {
	meta, ok := s.entitiesMeta[entity]
	if !ok {
		return
	}
	positions := s.getPositions(meta.aabb)
	for _, position := range positions {
		set, ok := s.positionsEntities[position]
		if !ok {
			continue
		}
		set.RemoveElements(entity)
		if len(set.Get()) == 0 {
			delete(s.positionsEntities, position)
		}
	}
	delete(s.entitiesMeta, entity)
}

func (s *orthoOcclusion) GetColliding(aabb transform.AABB) datastructures.Set[ecs.EntityID] {
	entities := datastructures.NewSet[ecs.EntityID]()
	positions := s.getPositions(aabb)
	for _, position := range positions {
		set, ok := s.positionsEntities[position]
		if !ok {
			continue
		}
		elements := set.Get()
		entities.Add(elements...)
	}
	return entities
}

type OcclusionSystem struct {
	ortho *orthoOcclusion

	// orthoEntities              datastructures.Set[ecs.EntityID]
	perspectiveEntities        datastructures.Set[ecs.EntityID]
	visibleOrthoEntities       datastructures.Set[ecs.EntityID]
	visiblePerspectiveEntities datastructures.Set[ecs.EntityID]
}

func NewOcclusionSystem(world ecs.World) *OcclusionSystem {
	// {
	// 	q := world.QueryEntitiesWithComponents(ecs.GetComponentType(transform.Transform{}))
	// 	makeVisible := func(ei []ecs.EntityID) {
	// 		for _, entity := range ei {
	// 			world.SaveComponent(entity, projection.Visible{})
	// 		}
	// 	}
	// 	q.OnAdd(makeVisible)
	// 	return nil
	// }

	// orthoEntities := datastructures.NewSet[ecs.EntityID]()
	perspectiveEntities := datastructures.NewSet[ecs.EntityID]()

	// visible ortho
	visibleOrthoEntities := datastructures.NewSet[ecs.EntityID]()
	visiblePerspectiveEntities := datastructures.NewSet[ecs.EntityID]()

	s := &OcclusionSystem{
		ortho: newOrthoOcclusion(mgl32.Vec2{1000, 1000}),
		// orthoEntities:              orthoEntities,
		perspectiveEntities:        perspectiveEntities,
		visibleOrthoEntities:       visibleOrthoEntities,
		visiblePerspectiveEntities: visiblePerspectiveEntities,
	}

	orthoCameraQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(projection.Ortho{}),
		ecs.GetComponentType(transform.Transform{}),
	)

	refresh := func() {
		cameraEntities := orthoCameraQuery.Entities()
		if len(cameraEntities) != 1 {
			return
		}
		camera := cameraEntities[0]
		cameraTransform, err := ecs.GetComponent[transform.Transform](world, camera)
		if err != nil {
			return
		}
		ortho, err := ecs.GetComponent[projection.Ortho](world, camera)
		if err != nil {
			return
		}
		cameraAabb := transform.NewAABB(
			// Min: mgl32.Vec3{left, bottom, -far},
			// Max: mgl32.Vec3{right, top, -near},
			mgl32.Vec3{cameraTransform.Pos.X() - ortho.Width/2, cameraTransform.Pos.Y() - ortho.Height/2, -ortho.Far},
			mgl32.Vec3{cameraTransform.Pos.X() + ortho.Width/2, cameraTransform.Pos.Y() + ortho.Height/2, -ortho.Near},
		)
		newVisibleOrthoEntities := s.ortho.GetColliding(cameraAabb)

		for _, entity := range visibleOrthoEntities.Get() {
			if _, ok := newVisibleOrthoEntities.GetIndex(entity); !ok {
				world.RemoveComponent(entity, ecs.GetComponentType(projection.Visible{}))
				visibleOrthoEntities.RemoveElements(entity)
			}
		}
		for _, entity := range newVisibleOrthoEntities.Get() {
			if _, ok := visibleOrthoEntities.GetIndex(entity); !ok {
				world.SaveComponent(entity, projection.Visible{})
				visibleOrthoEntities.Add(entity)
			}
		}
	}

	orthoCameraQuery.OnAdd(func(ei []ecs.EntityID) { refresh() })
	orthoCameraQuery.OnChange(func(ei []ecs.EntityID) { refresh() })
	orthoCameraQuery.OnRemove(func(ei []ecs.EntityID) { refresh() })

	entitiesQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(projection.UsedProjection{}),
		ecs.GetComponentType(transform.Transform{}),
	)
	entitiesQuery.OnAdd(func(entities []ecs.EntityID) {
		for _, entity := range entities {
			usedProj, err := ecs.GetComponent[projection.UsedProjection](world, entity)
			if err != nil {
				continue
			}
			transform, err := ecs.GetComponent[transform.Transform](world, entity)
			if err != nil {
				continue
			}
			if usedProj != projection.NewUsedProjection[projection.Ortho]() {
				world.SaveComponent(entity, projection.Visible{})
				continue
			}
			aabb := transform.ToAABB()
			s.ortho.Add(entity, aabb)
		}
		refresh()
	})
	entitiesQuery.OnChange(func(entities []ecs.EntityID) {
		for _, entity := range entities {
			usedProj, err := ecs.GetComponent[projection.UsedProjection](world, entity)
			if err != nil {
				continue
			}
			transform, err := ecs.GetComponent[transform.Transform](world, entity)
			if err != nil {
				continue
			}
			if usedProj != projection.NewUsedProjection[projection.Ortho]() {
				continue
			}
			aabb := transform.ToAABB()
			s.ortho.Remove(entity)
			s.ortho.Add(entity, aabb)
		}
		refresh()
	})
	entitiesQuery.OnRemove(func(entities []ecs.EntityID) {
		for _, entity := range entities {
			s.ortho.Remove(entity)
			world.RemoveComponent(entity, ecs.GetComponentType(projection.Visible{}))
		}
		refresh()
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

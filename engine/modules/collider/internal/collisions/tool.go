package collisions

import (
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"errors"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type tool struct {
	// shared
	logger        logger.Logger
	world         ecs.World
	transform     transform.Interface
	colliderArray ecs.ComponentsArray[collider.ColliderComponent]

	// detection
	groupsArray ecs.ComponentsArray[groups.GroupsComponent]
	assets      assets.Assets

	// tracking
	dirtySet          ecs.DirtySet
	chunkSize         float32
	chunks            map[mgl32.Vec2]datastructures.Set[ecs.EntityID]
	entitiesPositions map[ecs.EntityID][]mgl32.Vec2
}

func NewToolFactory(
	logger logger.Logger,
	assets assets.Assets,
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	chunkSize float32,
) ecs.ToolFactory[collider.ColliderTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(world ecs.World) collider.ColliderTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](world); ok {
			return t
		}
		dirtySet := ecs.NewDirtySet()
		transform := transformToolFactory.Build(world).Transform()
		transform.AddDirtySet(dirtySet)
		ecs.GetComponentsArray[collider.ColliderComponent](world).AddDirtySet(dirtySet)
		t := tool{
			logger:        logger,
			world:         world,
			transform:     transform,
			colliderArray: ecs.GetComponentsArray[collider.ColliderComponent](world),

			groupsArray: ecs.GetComponentsArray[groups.GroupsComponent](world),
			assets:      assets,

			dirtySet:          dirtySet,
			chunkSize:         chunkSize,
			chunks:            make(map[mgl32.Vec2]datastructures.Set[ecs.EntityID]),
			entitiesPositions: make(map[ecs.EntityID][]mgl32.Vec2),
		}
		world.SaveGlobal(t)
		return t
	})
}

func floorF32ToInt(num float32) int {
	// return int(math.Floor(float64(num)))
	return int(num)
}

func (c tool) getPositions(aabb collider.AABB) []mgl32.Vec2 {
	minPos, maxPos := aabb.Min.Vec2(), aabb.Max.Vec2()
	minGridX := floorF32ToInt(minPos.X() / c.chunkSize)
	minGridY := floorF32ToInt(minPos.Y() / c.chunkSize)
	maxGridX := floorF32ToInt(maxPos.X() / c.chunkSize)
	maxGridY := floorF32ToInt(maxPos.Y() / c.chunkSize)

	var positions []mgl32.Vec2
	for x := minGridX; x <= maxGridX; x++ {
		for y := minGridY; y <= maxGridY; y++ {
			tileX := float32(x) * c.chunkSize
			tileY := float32(y) * c.chunkSize
			tileCenter := mgl32.Vec2{tileX, tileY}

			positions = append(positions, tileCenter)
		}
	}

	return positions
}

func (t tool) ChunkSize() float32                                      { return t.chunkSize }
func (t tool) Chunks() map[mgl32.Vec2]datastructures.Set[ecs.EntityID] { return t.chunks }

// tracking

func (t tool) ApplyChanges() {
	entities := t.dirtySet.Get()
	t.Remove(entities...)
	for _, entity := range entities {
		if _, ok := t.colliderArray.Get(entity); !ok {
			continue
		}
		aabb := TransformAABB(t.transform, entity)
		positions := t.getPositions(aabb)
		t.entitiesPositions[entity] = positions
		for _, position := range positions {
			arr, ok := t.chunks[position]
			if !ok {
				arr = datastructures.NewSet[ecs.EntityID]()
			}
			arr.Add(entity)
			t.chunks[position] = arr
		}
	}
}

func (t tool) Remove(entities ...ecs.EntityID) {
	for _, entity := range entities {
		positions, ok := t.entitiesPositions[entity]
		if !ok {
			continue
		}
		delete(t.entitiesPositions, entity)
		for _, position := range positions {
			arr, ok := t.chunks[position]
			if !ok {
				continue
			}
			arr.RemoveElements(entity)
			if len(arr.Get()) == 0 {
				delete(t.chunks, position)
				continue
			}
		}
	}
}

//

func (t tool) Collider() collider.Interface { return t }

func (t tool) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (collider.ObjectRayCollision, error) {
	t.ApplyChanges()
	entityGroups, ok := t.groupsArray.Get(entity)
	if !ok {
		entityGroups = groups.DefaultGroups()
	}
	if entityGroups.GetSharedWith(ray.Groups).Mask == 0 {
		return nil, nil
	}

	aabb := TransformAABB(t.transform, entity)
	if ok, _ := RayAABBIntersect(ray, aabb); !ok {
		return nil, nil
	}

	colliderComponent, ok := t.colliderArray.Get(entity)
	if !ok {
		return nil, nil
	}
	colliderAsset, err := assets.GetAsset[collider.ColliderAsset](t.assets, colliderComponent.ID)
	if err != nil {
		return nil, err
	}

	//

	ray.Apply(t.transform.Mat4(entity).Inv())

	aabbs := colliderAsset.AABBs()
	ranges := colliderAsset.Ranges()
	polygons := colliderAsset.Polygons()

	rangesToVisit := []collider.Range{}
	if len(ranges) > 0 {
		rangesToVisit = append(rangesToVisit, collider.NewRange(collider.Branch, 0, 1))
	}

	var closestHit *collider.RayHit

	for len(rangesToVisit) > 0 {
		currentRange := rangesToVisit[len(rangesToVisit)-1]
		rangesToVisit = rangesToVisit[:len(rangesToVisit)-1]

		if currentRange.Target == collider.Branch {
			for i := currentRange.First; i < currentRange.First+currentRange.Count; i++ {
				aabb := aabbs[i]
				intersects, _ := RayAABBIntersect(ray, aabb)
				if !intersects {
					continue
				}
				rangesToVisit = append(rangesToVisit, ranges[i])
			}
		} else if currentRange.Target == collider.Leaf {
			polygons := polygons[currentRange.First : currentRange.First+currentRange.Count]
			for _, polygon := range polygons {
				intersect, dist := RayTriangleIntersect(ray, polygon)
				if !intersect {
					continue
				}
				if closestHit != nil && closestHit.Distance < dist {
					continue
				}

				ray := ray
				ray.MaxDistance = dist

				normal := polygon.B.Sub(polygon.A).Cross(polygon.C.Sub(polygon.A)).Normalize()
				hit := collider.NewRayHit(ray, normal)
				closestHit = &hit
			}
		}
	}

	if closestHit == nil {
		return nil, nil
	}

	collision := collider.NewObjectRayCollision(entity, *closestHit)
	return collision, nil
}

func (t tool) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) (collider.ObjectObjectCollision, error) {
	t.ApplyChanges()
	return nil, errors.New("501")
}

func (t tool) ShootRay(ray collider.Ray) (collider.ObjectRayCollision, error) {
	t.ApplyChanges()
	chunkSize := t.ChunkSize()
	gridX := floorF32ToInt(ray.Pos[0] / chunkSize)
	gridY := floorF32ToInt(ray.Pos[1] / chunkSize)
	chunkCoord := mgl32.Vec2{
		float32(gridX) * chunkSize,
		float32(gridY) * chunkSize,
	}

	chunk, ok := t.Chunks()[chunkCoord]
	if !ok {
		return nil, nil
	}

	var closestHit *collider.RayHit
	var closestEntity ecs.EntityID

	for _, entity := range chunk.Get() {
		collision, err := t.CollidesWithRay(entity, ray)
		if err != nil {
			return nil, err
		}
		if collision == nil {
			continue
		}

		if closestHit != nil && closestHit.Distance < collision.Hit().Distance {
			continue
		}

		hit := collision.Hit()
		closestHit = &hit
		closestEntity = entity
	}

	if closestHit == nil {
		return nil, nil
	}

	collision := collider.NewObjectRayCollision(closestEntity, *closestHit)
	return collision, nil
}

func (t tool) NarrowCollisions(entity ecs.EntityID) ([]ecs.EntityID, error) {
	t.ApplyChanges()
	return nil, errors.New("501")
}

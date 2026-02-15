package collisions

import (
	"engine/modules/assets"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"errors"
	"slices"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type service struct {
	// shared
	Logger    logger.Logger     `inject:"1"`
	World     ecs.World         `inject:"1"`
	Groups    groups.Service    `inject:"1"`
	Transform transform.Service `inject:"1"`
	Assets    assets.Service    `inject:"1"`

	colliderArray ecs.ComponentsArray[collider.Component]

	// tracking
	dirtySet          ecs.DirtySet
	chunkSize         float32
	chunks            map[mgl32.Vec2]datastructures.Set[ecs.EntityID]
	entitiesPositions map[ecs.EntityID][]mgl32.Vec2

	rayFallTroughPolicies []collider.FallTroughPolicy
}

func NewService(c ioc.Dic,
	chunkSize float32,
) collider.Service {
	t := ioc.GetServices[*service](c)

	t.dirtySet = ecs.NewDirtySet()
	t.colliderArray = ecs.GetComponentsArray[collider.Component](t.World)
	t.chunkSize = chunkSize
	t.chunks = make(map[mgl32.Vec2]datastructures.Set[ecs.EntityID])
	t.entitiesPositions = make(map[ecs.EntityID][]mgl32.Vec2)
	t.rayFallTroughPolicies = make([]collider.FallTroughPolicy, 0)

	t.Transform.AddDirtySet(t.dirtySet)
	colliderArray := ecs.GetComponentsArray[collider.Component](t.World)
	colliderArray.AddDirtySet(t.dirtySet)

	return t
}

func floorF32ToInt(num float32) int {
	// return int(math.Floor(float64(num)))
	return int(num)
}

func (c *service) getPositions(aabb collider.AABB) []mgl32.Vec2 {
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

func (t *service) ChunkSize() float32                                      { return t.chunkSize }
func (t *service) Chunks() map[mgl32.Vec2]datastructures.Set[ecs.EntityID] { return t.chunks }

// tracking

func (t *service) ApplyChanges() {
	entities := t.dirtySet.Get()
	t.Remove(entities...)
	for _, entity := range entities {
		if _, ok := t.colliderArray.Get(entity); !ok {
			continue
		}
		aabb := TransformAABB(t.Transform, entity)
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

func (t *service) Remove(entities ...ecs.EntityID) {
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

func (t *service) Component() ecs.ComponentsArray[collider.Component] { return t.colliderArray }

func (t *service) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) *collider.ObjectRayCollision {
	t.ApplyChanges()
	entityGroups, ok := t.Groups.Component().Get(entity)
	if !ok {
		entityGroups = groups.DefaultGroups()
	}
	if entityGroups.GetSharedWith(ray.Groups).Mask == 0 {
		return nil
	}

	aabb := TransformAABB(t.Transform, entity)
	if ok, _ := RayAABBIntersect(ray, aabb); !ok {
		return nil
	}

	colliderComponent, ok := t.colliderArray.Get(entity)
	if !ok {
		return nil
	}
	colliderAsset, err := assets.GetAsset[collider.ColliderAsset](t.Assets, colliderComponent.ID)
	if err != nil {
		// invalid internal state
		t.Logger.Warn(err)
		return nil
	}

	//

	ray.Apply(t.Transform.Mat4(entity).Inv())

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
		return nil
	}

	collision := collider.NewObjectRayCollision(entity, *closestHit)

	for _, rayFallTroughPolicy := range t.rayFallTroughPolicies {
		if fallThrough := rayFallTroughPolicy.FallThrough(collision); fallThrough {
			return nil
		}
	}
	return &collision
}

func (t *service) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) *collider.ObjectObjectCollision {
	t.ApplyChanges()
	t.Logger.Warn(errors.New("501"))
	return nil
}

func (t *service) Raycast(ray collider.Ray) *collider.ObjectRayCollision {
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
		return nil
	}

	var closestHit *collider.RayHit
	var closestEntity ecs.EntityID

	for _, entity := range chunk.Get() {
		collision := t.CollidesWithRay(entity, ray)
		if collision == nil {
			continue
		}

		if closestHit != nil && closestHit.Distance < collision.Hit.Distance {
			continue
		}

		hit := collision.Hit

		closestHit = &hit
		closestEntity = entity
	}

	if closestHit == nil {
		return nil
	}

	collision := collider.NewObjectRayCollision(closestEntity, *closestHit)
	return &collision
}

func (t *service) RaycastAll(ray collider.Ray) []collider.ObjectRayCollision {
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
		return nil
	}

	collisions := []collider.ObjectRayCollision{}

	for _, entity := range chunk.Get() {
		collision := t.CollidesWithRay(entity, ray)
		if collision == nil {
			continue
		}

		collisions = append(collisions, collider.NewObjectRayCollision(entity, collision.Hit))
	}

	slices.SortFunc(collisions, func(a, b collider.ObjectRayCollision) int {
		if a.Hit.Distance < b.Hit.Distance {
			return -1
		}
		if a.Hit.Distance > b.Hit.Distance {
			return 1
		}
		return 0
	})

	return collisions
}

func (t *service) NarrowCollisions(entity ecs.EntityID) []ecs.EntityID {
	t.ApplyChanges()
	t.Logger.Warn(errors.New("501"))
	return nil
}
func (t *service) AddRayFallThroughPolicy(rayFallTroughPolicy collider.FallTroughPolicy) {
	t.rayFallTroughPolicies = append(t.rayFallTroughPolicies, rayFallTroughPolicy)
}

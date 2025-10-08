package broadcollision

import (
	"errors"
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"math"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO 1
type CollisionDetectionService interface {
	// todo add collision groups
	// narrow
	CollidesWithRay(ecs.EntityID, collider.Ray) (ObjectRayCollision, error)
	CollidesWithObject(ecs.EntityID, ecs.EntityID) (ObjectObjectCollision, error)

	// broad
	ShootRay(collider.Ray) (ObjectRayCollision, error)
	NarrowCollisions(ecs.EntityID) ([]ecs.EntityID, error)
}

type collisionDetectionService struct {
	// world           ecs.World
	transformsArray ecs.ComponentsArray[transform.Transform]
	groupsArray     ecs.ComponentsArray[groups.Groups]
	colliderArray   ecs.ComponentsArray[collider.Collider]
	assets          assets.Assets
	worldCollider   worldCollider
	logger          logger.Logger
}

func newCollisionDetectionService(world ecs.World, assets assets.Assets, worldCollider worldCollider, logger logger.Logger) CollisionDetectionService {
	return &collisionDetectionService{
		// world,
		ecs.GetComponentsArray[transform.Transform](world.Components()),
		ecs.GetComponentsArray[groups.Groups](world.Components()),
		ecs.GetComponentsArray[collider.Collider](world.Components()),
		assets,
		worldCollider,
		logger,
	}
}

func (c *collisionDetectionService) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (ObjectRayCollision, error) {
	entityGroups, err := c.groupsArray.GetComponent(entity)
	if err != nil {
		entityGroups = groups.DefaultGroups()
	}
	if entityGroups.GetSharedWith(ray.Groups).Mask == 0 {
		return nil, nil
	}

	transformComponent, err := c.transformsArray.GetComponent(entity)
	if err != nil {
		return nil, err
	}
	if ok, _ := rayAABBIntersect(ray, collider.TransformAABB(transformComponent)); !ok {
		return nil, nil
	}

	colliderComponent, err := c.colliderArray.GetComponent(entity)
	if err != nil {
		return nil, err
	}
	colliderAsset, err := assets.GetAsset[collider.ColliderAsset](c.assets, colliderComponent.ID)
	if err != nil {
		return nil, err
	}

	//

	ray.Apply(transformComponent.Mat4().Inv())

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
				intersects, _ := rayAABBIntersect(ray, aabb)
				if !intersects {
					continue
				}
				rangesToVisit = append(rangesToVisit, ranges[i])
			}
		} else if currentRange.Target == collider.Leaf {
			polygons := polygons[currentRange.First : currentRange.First+currentRange.Count]
			for _, polygon := range polygons {
				intersect, dist := rayTriangleIntersect(ray, polygon)
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

	collision := NewObjectRayCollision(entity, *closestHit)
	return collision, nil
}

func (c *collisionDetectionService) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) (ObjectObjectCollision, error) {
	return nil, errors.New("501")
}

func (c *collisionDetectionService) ShootRay(ray collider.Ray) (ObjectRayCollision, error) {
	chunkCoordBefore := ray.Pos.Mul(1 / c.worldCollider.ChunkSize())
	chunkCoord := mgl32.Vec2{
		float32(math.Round(float64(chunkCoordBefore[0]))) * c.worldCollider.ChunkSize(),
		float32(math.Round(float64(chunkCoordBefore[1]))) * c.worldCollider.ChunkSize(),
	}

	chunk, ok := c.worldCollider.Chunks()[chunkCoord]
	if !ok {
		return nil, nil
	}

	var closestHit *collider.RayHit
	var closestEntity ecs.EntityID

	for _, entity := range chunk.Get() {
		collision, err := c.CollidesWithRay(entity, ray)
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

	collision := NewObjectRayCollision(closestEntity, *closestHit)
	return collision, nil
}

func (c *collisionDetectionService) NarrowCollisions(entity ecs.EntityID) ([]ecs.EntityID, error) {
	return nil, errors.New("501")
}

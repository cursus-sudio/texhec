package collisions

import (
	"errors"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"frontend/modules/transform"
	"frontend/services/assets"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type collisionDetectionService struct {
	world                ecs.World
	transformTransaction transform.TransformTransaction
	groupsArray          ecs.ComponentsArray[groups.GroupsComponent]
	colliderArray        ecs.ComponentsArray[collider.ColliderComponent]
	assets               assets.Assets
	worldCollider        worldCollider
	logger               logger.Logger
}

func newCollisionDetectionService(
	world ecs.World,
	transformTransaction transform.TransformTransaction,
	assets assets.Assets,
	worldCollider worldCollider,
	logger logger.Logger,
) collider.CollisionTool {
	return &collisionDetectionService{
		world,
		transformTransaction,
		ecs.GetComponentsArray[groups.GroupsComponent](world.Components()),
		ecs.GetComponentsArray[collider.ColliderComponent](world.Components()),
		assets,
		worldCollider,
		logger,
	}
}

func (c *collisionDetectionService) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (collider.ObjectRayCollision, error) {
	entityGroups, err := c.groupsArray.GetComponent(entity)
	if err != nil {
		entityGroups = groups.DefaultGroups()
	}
	if entityGroups.GetSharedWith(ray.Groups).Mask == 0 {
		return nil, nil
	}

	transform := c.transformTransaction.GetEntity(entity)
	aabb, err := TransformAABB(transform)
	if err != nil {
		return nil, err
	}
	if ok, _ := RayAABBIntersect(ray, aabb); !ok {
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

	ray.Apply(transform.Mat4().Inv())

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

func (c *collisionDetectionService) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) (collider.ObjectObjectCollision, error) {
	return nil, errors.New("501")
}

func (c *collisionDetectionService) ShootRay(ray collider.Ray) (collider.ObjectRayCollision, error) {
	chunkSize := c.worldCollider.ChunkSize()
	gridX := floorF32ToInt(ray.Pos[0] / chunkSize)
	gridY := floorF32ToInt(ray.Pos[1] / chunkSize)
	chunkCoord := mgl32.Vec2{
		float32(gridX) * chunkSize,
		float32(gridY) * chunkSize,
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

	collision := collider.NewObjectRayCollision(closestEntity, *closestHit)
	return collision, nil
}

func (c *collisionDetectionService) NarrowCollisions(entity ecs.EntityID) ([]ecs.EntityID, error) {
	return nil, errors.New("501")
}

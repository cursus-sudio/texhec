package broadcollision

import (
	"errors"
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"math"
)

// TODO 1
type CollisionDetectionService interface {
	// narrow
	CollidesWithRay(ecs.EntityID, collider.Ray) (ObjectRayCollision, error)
	CollidesWithObject(ecs.EntityID, ecs.EntityID) (ObjectObjectCollision, error)

	// broad
	ShootRay(collider.Ray) (ObjectRayCollision, error)
	NarrowCollisions(ecs.EntityID) ([]ecs.EntityID, error)
}

type collisionDetectionService struct {
	world         ecs.World
	assets        assets.Assets
	worldCollider worldCollider
}

func newCollisionDetectionService(world ecs.World, assets assets.Assets, worldCollider worldCollider) CollisionDetectionService {
	return &collisionDetectionService{world, assets, worldCollider}
}

func (c *collisionDetectionService) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (ObjectRayCollision, error) {
	transformComponent, err := ecs.GetComponent[transform.Transform](c.world, entity)
	if err != nil {
		return nil, err
	}
	if ok, _ := rayAABBIntersect(
		ray, collider.TransformAABB(transformComponent), math.MaxFloat32); !ok {
		return nil, nil
	}

	colliderComponent, err := ecs.GetComponent[collider.Collider](c.world, entity)
	if err != nil {
		return nil, err
	}
	colliderAsset, err := assets.GetAsset[collider.ColliderAsset](c.assets, colliderComponent.ID)
	if err != nil {
		return nil, err
	}

	//

	aabbs := colliderAsset.AABBs()
	ranges := colliderAsset.Ranges()
	polygons := colliderAsset.Polygons()

	rangesToVisit := []collider.Range{}
	if len(ranges) > 0 {
		rangesToVisit = append(rangesToVisit, ranges[0])
	}

	var closestHit *collider.RayHit
	var minDistance float32 = math.MaxFloat32

	for len(rangesToVisit) > 0 {
		currentRange := rangesToVisit[len(rangesToVisit)-1]
		rangesToVisit = rangesToVisit[:len(rangesToVisit)-1]

		aabb := aabbs[currentRange.First]

		intersects, _ := rayAABBIntersect(ray, aabb, minDistance)
		if !intersects {
			continue
		}

		if currentRange.Target == collider.Branch {
			for i := uint32(0); i < currentRange.Count; i++ {
				childIndex := currentRange.First + 1 + i
				if childIndex < uint32(len(ranges)) {
					rangesToVisit = append(rangesToVisit, ranges[childIndex])
				}
			}
			continue
		}

		for i := uint32(0); i < currentRange.Count; i++ {
			polygonIndex := currentRange.First + i
			polygon := polygons[polygonIndex]

			intersect, dist := rayTriangleIntersect(ray, polygon, minDistance)
			if !intersect {
				continue
			}

			if dist >= minDistance {
				continue
			}

			minDistance = dist
			hitPoint := ray.Pos.Add(ray.Direction.Mul(dist))
			normal := polygon.B.Sub(polygon.A).Cross(polygon.C.Sub(polygon.A)).Normalize()

			hit := collider.NewRayHit(hitPoint, normal, dist)
			closestHit = &hit
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
	aabbs := c.worldCollider.AABBs()
	ranges := c.worldCollider.Ranges()
	entities := c.worldCollider.Entities()
	if len(aabbs) == 0 {
		return nil, nil
	}

	branchesToVisit := []int{}
	branchesToVisit = append(branchesToVisit, 0)

	var closestHit *collider.RayHit
	var closestEntity ecs.EntityID
	var minDistance float32 = math.MaxFloat32

	for len(branchesToVisit) > 0 {
		currentBranch := branchesToVisit[len(branchesToVisit)-1]
		branchesToVisit = branchesToVisit[:len(branchesToVisit)-1]
		currentRange := ranges[currentBranch]

		intersects, dist := rayAABBIntersect(ray, aabbs[currentRange.First], minDistance)
		if !intersects || dist >= closestHit.Distance {
			continue
		}

		if currentRange.Target == collider.Branch {
			for i := uint32(0); i < currentRange.Count; i++ {
				childIndex := currentRange.First + i
				branchesToVisit = append(branchesToVisit, int(childIndex))
			}
			continue
		}

		// leaf
		for i := uint32(0); i < currentRange.Count; i++ {
			entityIndex := currentRange.First + i
			entity := entities[entityIndex]

			collision, err := c.CollidesWithRay(entity, ray)
			if err != nil {
				return nil, err
			}
			if collision == nil {
				continue
			}

			if collision.Hit().Distance >= minDistance {
				continue
			}

			minDistance = dist
			hit := collision.Hit()
			closestHit = &hit
			closestEntity = entity
		}
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

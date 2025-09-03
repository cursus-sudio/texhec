package broadcollision

import (
	"errors"
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"math"

	"github.com/go-gl/mathgl/mgl32"
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
	if ok, _ := rayAABBIntersect(ray, collider.TransformAABB(transformComponent), math.MaxFloat32); !ok {
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

	{
		inverseTransform := transformComponent.Mat4().Inv()
		inverseTranspose := inverseTransform.Transpose()
		ray = collider.NewRay(
			inverseTransform.Mul4x1(ray.Pos.Vec4(1.0)).Vec3(),
			inverseTranspose.Mul4x1(ray.Direction.Vec4(1.0)).Vec3(),
		)
	}

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
				var maxDistance float32 = math.MaxFloat32
				if closestHit != nil {
					maxDistance = closestHit.Distance
				}
				intersects, _ := rayAABBIntersect(ray, aabb, maxDistance)
				if !intersects {
					continue
				}
				rangesToVisit = append(rangesToVisit, ranges[i])
			}
		} else if currentRange.Target == collider.Leaf {
			polygons := polygons[currentRange.First : currentRange.First+currentRange.Count]
			for _, polygon := range polygons {
				var maxDistance float32 = math.MaxFloat32
				if closestHit != nil {
					maxDistance = closestHit.Distance
				}
				intersect, dist := rayTriangleIntersect(ray, polygon, maxDistance)
				if !intersect {
					continue
				}

				if closestHit != nil && dist >= closestHit.Distance {
					continue
				}

				hitPoint := ray.Pos.Add(ray.Direction.Mul(dist))
				normal := polygon.B.Sub(polygon.A).Cross(polygon.C.Sub(polygon.A)).Normalize()

				hit := collider.NewRayHit(hitPoint, normal, dist)
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

		if closestHit != nil && collision.Hit().Distance <= closestHit.Distance {
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

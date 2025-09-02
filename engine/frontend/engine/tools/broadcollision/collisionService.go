package broadcollision

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
)

type CollisionServiceFactory func(ecs.World) CollisionService

type CollisionService interface {
	CollisionDetectionService
	CollidersTrackingService
}

// TODO 3
// change visibility check in projection/occlusion

type collisionsService struct {
	world                ecs.World
	assets               assets.Assets
	staticWorldCollider  worldCollider
	dynamicWorldCollider worldCollider
}

var factory CollisionServiceFactory = func(w ecs.World) CollisionService {
	return &collisionsService{
		world:                w,
		staticWorldCollider:  newWorldCollider(w),
		dynamicWorldCollider: newWorldCollider(w),
	}
}

func (s *collisionsService) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (ObjectRayCollision, error) {
	return newCollisionDetectionService(s.world, s.assets, nil).
		CollidesWithRay(entity, ray)
}
func (s *collisionsService) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) (ObjectObjectCollision, error) {
	return newCollisionDetectionService(s.world, s.assets, nil).
		CollidesWithObject(entityA, entityB)
}

func (s *collisionsService) ShootRay(entity collider.Ray) (ObjectRayCollision, error) {
	c1, err := newCollisionDetectionService(s.world, s.assets, s.staticWorldCollider).ShootRay(entity)
	if c1 == nil || err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.assets, s.dynamicWorldCollider).ShootRay(entity)
	if c2 == nil || err != nil {
		return nil, err
	}
	if c1.Hit().Distance < c2.Hit().Distance {
		return c1, nil
	}
	return c2, nil
}
func (s *collisionsService) NarrowCollisions(entity ecs.EntityID) ([]ecs.EntityID, error) {
	c1, err := newCollisionDetectionService(s.world, s.assets, s.staticWorldCollider).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.assets, s.dynamicWorldCollider).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	collisions := append(c1, c2...)
	return collisions, nil
}

func (s *collisionsService) Add(entities ...ecs.EntityID) {
	static := make([]ecs.EntityID, 0, len(entities))
	dynamic := make([]ecs.EntityID, 0, len(entities))
	for _, entity := range entities {
		_, err := ecs.GetComponent[transform.Static](s.world, entity)
		isDynamic := err != nil
		if isDynamic {
			dynamic = append(dynamic, entity)
		} else {
			static = append(static, entity)
		}
	}
	s.dynamicWorldCollider.Add(dynamic...)
	s.staticWorldCollider.Add(static...)

}
func (s *collisionsService) Update(entities ...ecs.EntityID) {
	s.dynamicWorldCollider.Update(entities...)
}
func (s *collisionsService) Remove(entities ...ecs.EntityID) {
	s.dynamicWorldCollider.Remove(entities...)
}

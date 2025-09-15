package broadcollision

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"shared/services/ecs"
	"shared/services/logger"
	"sync"
)

type CollisionServiceFactory func(ecs.World) CollisionService

type CollisionService interface {
	CollisionDetectionService
	CollidersTrackingService
}

// TODO 3
// change visibility check in projection/occlusion

type register struct {
	staticWorldCollider  worldCollider
	dynamicWorldCollider worldCollider
	mutex                sync.Locker
}

func newRegister(world ecs.World) register {
	r := register{
		newWorldCollider(world, 100),
		newWorldCollider(world, 100),
		&sync.Mutex{},
	}
	world.SaveRegister(r)
	return r
}

type collisionsService struct {
	world                ecs.World
	transformStaticArray ecs.ComponentsArray[transform.Static]
	assets               assets.Assets
	logger               logger.Logger
}

func factory(assets assets.Assets, logger logger.Logger) CollisionServiceFactory {
	return func(w ecs.World) CollisionService {
		return &collisionsService{
			world:                w,
			transformStaticArray: ecs.GetComponentsArray[transform.Static](w.Components()),
			assets:               assets,
			logger:               logger,
		}
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

func (s *collisionsService) ShootRay(ray collider.Ray) (ObjectRayCollision, error) {
	r, err := ecs.GetRegister[register](s.world)
	if err != nil {
		r = newRegister(s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c1, err := newCollisionDetectionService(s.world, s.assets, r.staticWorldCollider).ShootRay(ray)
	if err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.assets, r.dynamicWorldCollider).ShootRay(ray)
	if err != nil {
		return nil, err
	}
	if c1 == nil {
		return c2, nil
	}
	if c2 == nil {
		return c1, nil
	}
	if c1.Hit().Distance < c2.Hit().Distance {
		return c1, nil
	}
	return c2, nil
}
func (s *collisionsService) NarrowCollisions(entity ecs.EntityID) ([]ecs.EntityID, error) {
	r, err := ecs.GetRegister[register](s.world)
	if err != nil {
		r = newRegister(s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c1, err := newCollisionDetectionService(s.world, s.assets, r.staticWorldCollider).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.assets, r.dynamicWorldCollider).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	collisions := append(c1, c2...)
	return collisions, nil
}

func (s *collisionsService) Add(entities ...ecs.EntityID) {
	r, err := ecs.GetRegister[register](s.world)
	if err != nil {
		r = newRegister(s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	static := make([]ecs.EntityID, 0, len(entities))
	dynamic := make([]ecs.EntityID, 0, len(entities))
	for _, entity := range entities {

		_, err := s.transformStaticArray.GetComponent(entity)
		isDynamic := err != nil
		if isDynamic {
			dynamic = append(dynamic, entity)
		} else {
			static = append(static, entity)
		}
	}
	r.dynamicWorldCollider.Add(dynamic...)
	r.staticWorldCollider.Add(static...)

}
func (s *collisionsService) Update(entities ...ecs.EntityID) {
	r, err := ecs.GetRegister[register](s.world)
	if err != nil {
		r = newRegister(s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.dynamicWorldCollider.Update(entities...)
}
func (s *collisionsService) Remove(entities ...ecs.EntityID) {
	r, err := ecs.GetRegister[register](s.world)
	if err != nil {
		r = newRegister(s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.dynamicWorldCollider.Remove(entities...)
}

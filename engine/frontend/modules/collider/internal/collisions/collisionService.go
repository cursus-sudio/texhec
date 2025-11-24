package collisions

import (
	"frontend/modules/collider"
	"frontend/modules/transform"
	"frontend/services/assets"
	"shared/services/ecs"
	"shared/services/logger"
	"sync"
)

type CollisionService interface {
	collider.CollisionTool
	CollidersTrackingService
}

type global struct {
	staticWorldCollider  worldCollider
	dynamicWorldCollider worldCollider
	mutex                sync.Locker
}

func newRegister(
	logger logger.Logger,
	transformTransaction transform.Transaction,
	world ecs.World,
) global {
	r := global{
		newWorldCollider(logger, world, transformTransaction, 100),
		newWorldCollider(logger, world, transformTransaction, 100),
		&sync.Mutex{},
	}
	world.SaveGlobal(r)
	return r
}

type collisionsService struct {
	world                ecs.World
	transformTransaction transform.Transaction
	transformStaticArray ecs.ComponentsArray[transform.StaticComponent]
	assets               assets.Assets
	logger               logger.Logger
}

func Factory(
	assets assets.Assets,
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.Tool],
) ecs.ToolFactory[CollisionService] {
	return ecs.NewToolFactory(func(w ecs.World) CollisionService {
		return &collisionsService{
			world:                w,
			transformTransaction: transformToolFactory.Build(w).Transaction(),
			transformStaticArray: ecs.GetComponentsArray[transform.StaticComponent](w),
			assets:               assets,
			logger:               logger,
		}
	})
}

func (s *collisionsService) CollidesWithRay(entity ecs.EntityID, ray collider.Ray) (collider.ObjectRayCollision, error) {
	return newCollisionDetectionService(s.world, s.transformTransaction, s.assets, nil, s.logger).
		CollidesWithRay(entity, ray)
}
func (s *collisionsService) CollidesWithObject(entityA ecs.EntityID, entityB ecs.EntityID) (collider.ObjectObjectCollision, error) {
	return newCollisionDetectionService(s.world, s.transformTransaction, s.assets, nil, s.logger).
		CollidesWithObject(entityA, entityB)
}

func (s *collisionsService) ShootRay(ray collider.Ray) (collider.ObjectRayCollision, error) {
	r, err := ecs.GetGlobal[global](s.world)
	if err != nil {
		r = newRegister(s.logger, s.transformTransaction, s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c1, err := newCollisionDetectionService(s.world, s.transformTransaction, s.assets, r.staticWorldCollider, s.logger).ShootRay(ray)
	if err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.transformTransaction, s.assets, r.dynamicWorldCollider, s.logger).ShootRay(ray)
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
	r, err := ecs.GetGlobal[global](s.world)
	if err != nil {
		r = newRegister(s.logger, s.transformTransaction, s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	c1, err := newCollisionDetectionService(s.world, s.transformTransaction, s.assets, r.staticWorldCollider, s.logger).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	c2, err := newCollisionDetectionService(s.world, s.transformTransaction, s.assets, r.dynamicWorldCollider, s.logger).NarrowCollisions(entity)
	if err != nil {
		return nil, err
	}
	collisions := append(c1, c2...)
	return collisions, nil
}

func (s *collisionsService) Add(entities ...ecs.EntityID) {
	r, err := ecs.GetGlobal[global](s.world)
	if err != nil {
		r = newRegister(s.logger, s.transformTransaction, s.world)
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
	r, err := ecs.GetGlobal[global](s.world)
	if err != nil {
		r = newRegister(s.logger, s.transformTransaction, s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.dynamicWorldCollider.Update(entities...)
}
func (s *collisionsService) Remove(entities ...ecs.EntityID) {
	r, err := ecs.GetGlobal[global](s.world)
	if err != nil {
		r = newRegister(s.logger, s.transformTransaction, s.world)
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.dynamicWorldCollider.Remove(entities...)
}

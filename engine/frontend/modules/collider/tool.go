package collider

import "shared/services/ecs"

type CollisionTool interface {
	// todo add collision groups
	// narrow
	CollidesWithRay(ecs.EntityID, Ray) (ObjectRayCollision, error)
	CollidesWithObject(ecs.EntityID, ecs.EntityID) (ObjectObjectCollision, error)

	// broad
	ShootRay(Ray) (ObjectRayCollision, error)
	NarrowCollisions(ecs.EntityID) ([]ecs.EntityID, error)
}

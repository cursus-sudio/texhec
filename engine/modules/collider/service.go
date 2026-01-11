package collider

import (
	"engine/services/ecs"
)

type Service interface {
	Component() ecs.ComponentsArray[Component]

	// todo add collision groups
	// narrow
	CollidesWithRay(ecs.EntityID, Ray) *ObjectRayCollision
	CollidesWithObject(ecs.EntityID, ecs.EntityID) *ObjectObjectCollision

	// broad
	Raycast(Ray) *ObjectRayCollision
	RaycastAll(Ray) []ObjectRayCollision
	NarrowCollisions(ecs.EntityID) []ecs.EntityID
}

// ```go
//     ShootRay(Ray) (ObjectRayCollision, error)
//     ShootRaycast(Ray) []ObjectRayCollision
// ```
//
// is this naming clear ?
// shootRay checks nearest collision returns error if there is none.
// shootRaycast returns all matching collisions if there is none returns empty slice.
//
// could shootRaycast be named better or this naming is already clear ?

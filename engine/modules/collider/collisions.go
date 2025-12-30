package collider

import (
	"engine/services/ecs"
)

type ObjectRayCollision struct {
	Entity ecs.EntityID
	Hit    RayHit
}

func NewObjectRayCollision(entity ecs.EntityID, hit RayHit) ObjectRayCollision {
	return ObjectRayCollision{entity, hit}
}

//

type ObjectObjectCollision struct {
	PolygonPairs [][2]Polygon
}

func NewObjectObjectCollision(pairs [][2]Polygon) ObjectObjectCollision {
	return ObjectObjectCollision{PolygonPairs: pairs}
}

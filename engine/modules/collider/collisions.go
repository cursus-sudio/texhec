package collider

import (
	"engine/services/ecs"
)

type ObjectRayCollision interface {
	Entity() ecs.EntityID
	Hit() RayHit
}

type objectRayCollision struct {
	entity ecs.EntityID
	hit    RayHit
}

func NewObjectRayCollision(entity ecs.EntityID, hit RayHit) ObjectRayCollision {
	return &objectRayCollision{entity, hit}
}

func (c *objectRayCollision) Entity() ecs.EntityID { return c.entity }
func (c *objectRayCollision) Hit() RayHit          { return c.hit }

//

type ObjectObjectCollision interface {
	PolygonPairs() [][2]Polygon
}

type objectObjectCollision struct {
	pairs [][2]Polygon
}

func NewObjectObjectCollision(pairs [][2]Polygon) ObjectObjectCollision {
	return &objectObjectCollision{pairs: pairs}
}

func (c *objectObjectCollision) PolygonPairs() [][2]Polygon { return c.pairs }

package broadcollision

import (
	"frontend/engine/components/collider"
	"shared/services/ecs"
)

type ObjectRayCollision interface {
	Entity() ecs.EntityID
	Hit() collider.RayHit
}

type objectRayCollision struct {
	entity ecs.EntityID
	hit    collider.RayHit
}

func NewObjectRayCollision(entity ecs.EntityID, hit collider.RayHit) ObjectRayCollision {
	return &objectRayCollision{entity, hit}
}

func (c *objectRayCollision) Entity() ecs.EntityID { return c.entity }
func (c *objectRayCollision) Hit() collider.RayHit { return c.hit }

//

type ObjectObjectCollision interface {
	PolygonPairs() [][2]collider.Polygon
}

type objectObjectCollision struct {
	pairs [][2]collider.Polygon
}

func NewObjectObjectCollision(pairs [][2]collider.Polygon) ObjectObjectCollision {
	return &objectObjectCollision{pairs: pairs}
}

func (c *objectObjectCollision) PolygonPairs() [][2]collider.Polygon { return c.pairs }

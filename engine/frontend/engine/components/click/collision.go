package click

import "frontend/services/ecs"

type Collider struct {
	// vertices ?
	// its aabb box ?
	// circle ?
}

type OnCollision struct {
	// store collisions not only entities
	// collision = 2 points nearest to the other shape center (at least in simple shapes)
	CollidingEntities []ecs.EntityId
	EmitEvents        []any
}

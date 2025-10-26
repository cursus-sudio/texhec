package camera

import "shared/services/ecs"

type Camera struct {
	Projection ecs.ComponentType
}

func NewCamera(projection ecs.ComponentType) Camera {
	return Camera{projection}
}

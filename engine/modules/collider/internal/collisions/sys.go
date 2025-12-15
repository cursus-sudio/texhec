package collisions

import (
	"engine/modules/collider"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
)

func NewColliderSystem(
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.Transform],
	serviceFactory ecs.ToolFactory[collider.Collider],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		serviceFactory.Build(w)
		return nil
	})
}

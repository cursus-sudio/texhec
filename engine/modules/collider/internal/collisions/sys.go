package collisions

import (
	"engine/modules/collider"
	"engine/services/ecs"
	"engine/services/logger"
)

func NewColliderSystem(
	logger logger.Logger,
	serviceFactory ecs.ToolFactory[collider.World, collider.ColliderTool],
) ecs.SystemRegister[collider.World] {
	return ecs.NewSystemRegister(func(w collider.World) error {
		serviceFactory.Build(w)
		return nil
	})
}

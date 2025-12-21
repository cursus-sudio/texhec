package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/ecs"
)

func newWorld(
	uuidToolFactory uuid.ToolFactory,
) record.World {
	type world struct {
		ecs.World
		uuid.UUIDTool
	}

	w := world{
		World: ecs.NewWorld(),
	}
	w.UUIDTool = uuidToolFactory.Build(w)
	return w
}

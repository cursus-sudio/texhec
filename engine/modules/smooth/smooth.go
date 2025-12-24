package smooth

import (
	"engine/modules/record"
	"engine/modules/transition"
	"engine/services/ecs"
)

type StartSystem ecs.SystemRegister[World]
type StopSystem ecs.SystemRegister[World]

type World interface {
	ecs.World
	transition.TransitionTool
	record.RecordTool
}

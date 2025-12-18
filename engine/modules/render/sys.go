package render

import "engine/services/ecs"

type System ecs.SystemRegister[World]

type RenderEvent struct{}

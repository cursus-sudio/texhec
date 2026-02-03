package render

import "engine/services/ecs"

type System ecs.SystemRegister
type SystemRenderer ecs.SystemRegister

type FlushEvent struct{}
type RenderEvent struct{}

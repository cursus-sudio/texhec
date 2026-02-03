package render

import "engine/services/ecs"

type System ecs.SystemRegister
type SystemRenderer ecs.SystemRegister

type RenderEvent struct {
	Camera ecs.EntityID
}

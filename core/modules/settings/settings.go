package settings

import (
	"core/modules/ui"
	"engine"
	"engine/services/ecs"
)

type World interface {
	engine.World
	ui.UiTool
}

type System ecs.SystemRegister[World]

type EnterSettingsEvent struct{}

type EnterSettingsForParentEvent struct {
	Parent ecs.EntityID
}

package settings

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister

type EnterSettingsEvent struct{}

type EnterSettingsForParentEvent struct {
	Parent ecs.EntityID
}

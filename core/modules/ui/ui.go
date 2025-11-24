package ui

import "shared/services/ecs"

type System ecs.SystemRegister

// marker which says module relative to which element to position
type UiCameraComponent struct{}

type UnselectEvent struct{}

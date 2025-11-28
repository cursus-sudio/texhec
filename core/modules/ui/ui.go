package ui

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister

// marker which says module relative to which element to position
type UiCameraComponent struct{}

type HideUiEvent struct{}

type Tool interface {
	// returns parent to attach ui elements
	// potentially with enter animation
	Show() (parent ecs.EntityID)
	// removes all children
	Hide()
}

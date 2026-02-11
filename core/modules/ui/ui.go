package ui

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister

// marker which says module relative to which element to position
type UiCameraComponent struct{}

type HideUiEvent struct{}

type Button struct {
	Text  string
	Event any
}

func NewButton(text string, event any) Button {
	return Button{text, event}
}

type Service interface {
	UiCamera() ecs.ComponentsArray[UiCameraComponent]
	AnimatedBackground() ecs.ComponentsArray[AnimatedBackgroundComponent]
	CursorCamera() ecs.ComponentsArray[CursorCameraComponent]
	// returns parent to attach ui elements
	// potentially with enter animation
	Show() (parents []ecs.EntityID)
	// removes all children
	Hide()

	// elements
	// Buttons(...Button) []ecs.EntityID
}

package ui

import (
	"engine"
	"engine/services/ecs"
)

type System ecs.SystemRegister[World]

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

type UiTool interface {
	Ui() Interface
}

type World interface {
	engine.World
}

type Interface interface {
	UiCamera() ecs.ComponentsArray[UiCameraComponent]
	// returns parent to attach ui elements
	// potentially with enter animation
	Show() (parent ecs.EntityID)
	// removes all children
	Hide()

	// elements
	// Buttons(...Button) []ecs.EntityID
}

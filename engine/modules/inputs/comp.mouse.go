package inputs

import "engine/services/ecs"

type HoveredComponent struct {
	Camera ecs.EntityID
}

func NewHovered(camera ecs.EntityID) HoveredComponent {
	return HoveredComponent{camera}
}

//

type DraggedComponent struct {
	Camera ecs.EntityID
}

func NewDragged(camera ecs.EntityID) DraggedComponent {
	return DraggedComponent{camera}
}

//

// keeps element selected even if user drags outside
type KeepSelectedComponent struct{}

//

type MouseLeftClickComponent struct{ Event any }
type MouseDoubleLeftClickComponent struct{ Event any }

type MouseRightClickComponent struct{ Event any }
type MouseDoubleRightClickComponent struct{ Event any }

type MouseEnterComponent struct{ Event any }
type MouseLeaveComponent struct{ Event any }

type MouseHoverComponent struct{ Event any }
type MouseDragComponent struct{ Event any }

func NewMouseLeftClick(e any) MouseLeftClickComponent { return MouseLeftClickComponent{e} }
func NewMouseDoubleLeftClick(e any) MouseDoubleLeftClickComponent {
	return MouseDoubleLeftClickComponent{e}
}

func NewMouseRightClick(e any) MouseRightClickComponent { return MouseRightClickComponent{e} }
func NewMouseDoubleRightClick(e any) MouseDoubleRightClickComponent {
	return MouseDoubleRightClickComponent{e}
}

func NewMouseEnterComponent(event any) MouseEnterComponent { return MouseEnterComponent{event} }
func NewMouseLeaveComponent(event any) MouseLeaveComponent { return MouseLeaveComponent{event} }

func NewMouseHoverComponent(event any) MouseHoverComponent { return MouseHoverComponent{event} }
func NewMouseDragComponent(event any) MouseDragComponent   { return MouseDragComponent{event} }

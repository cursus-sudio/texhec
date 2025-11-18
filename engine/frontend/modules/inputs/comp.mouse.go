package inputs

import "shared/services/ecs"

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

//

// keeps element selected even if user drags outside
type KeepSelectedComponent struct{}

//

type MouseEventsComponent struct {
	LeftClickEvents        []any
	DoubleLeftClickEvents  []any
	RightClickEvents       []any
	DoubleRightClickEvents []any

	MouseEnterEvents []any
	MouseLeaveEvents []any

	// hover event is triggered every frame object is hovered
	HoverEvents []any
	DragEvents  []any
}

func NewMouseEvents() MouseEventsComponent {
	return MouseEventsComponent{}
}

func (comp MouseEventsComponent) Ptr() *MouseEventsComponent { return &comp }
func (comp *MouseEventsComponent) Val() MouseEventsComponent { return *comp }

func (component *MouseEventsComponent) AddLeftClickEvents(events ...any) *MouseEventsComponent {
	component.LeftClickEvents = append(component.LeftClickEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddDoubleLeftClickEvents(events ...any) *MouseEventsComponent {
	component.DoubleLeftClickEvents = append(component.DoubleLeftClickEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddRightClickEvents(events ...any) *MouseEventsComponent {
	component.RightClickEvents = append(component.RightClickEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddDoubleRightClickEvents(events ...any) *MouseEventsComponent {
	component.DoubleRightClickEvents = append(component.DoubleRightClickEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddMouseEnterEvents(events ...any) *MouseEventsComponent {
	component.MouseEnterEvents = append(component.MouseEnterEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddMouseLeaveEvents(events ...any) *MouseEventsComponent {
	component.MouseLeaveEvents = append(component.MouseLeaveEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddHoverEvents(events ...any) *MouseEventsComponent {
	component.HoverEvents = append(component.HoverEvents, events...)
	return component
}

func (component *MouseEventsComponent) AddDragEvents(events ...any) *MouseEventsComponent {
	component.DragEvents = append(component.DragEvents, events...)
	return component
}

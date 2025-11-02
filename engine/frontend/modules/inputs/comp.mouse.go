package inputs

type HoveredComponent struct{}

func NewHovered() HoveredComponent { return HoveredComponent{} }

//

type DraggedComponent struct{}

func NewDragged() DraggedComponent { return DraggedComponent{} }

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

func (component MouseEventsComponent) Clone() MouseEventsComponent {
	return MouseEventsComponent{
		LeftClickEvents:        component.LeftClickEvents,
		DoubleLeftClickEvents:  component.DoubleLeftClickEvents,
		RightClickEvents:       component.RightClickEvents,
		DoubleRightClickEvents: component.DoubleRightClickEvents,

		MouseEnterEvents: component.MouseEnterEvents,
		MouseLeaveEvents: component.MouseLeaveEvents,

		HoverEvents: component.HoverEvents,
		DragEvents:  component.DragEvents,
	}
}

func (component MouseEventsComponent) AddLeftClickEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.LeftClickEvents = append(r.LeftClickEvents, events...)
	return r
}

func (component MouseEventsComponent) AddDoubleLeftClickEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.DoubleLeftClickEvents = append(r.DoubleLeftClickEvents, events...)
	return r
}

func (component MouseEventsComponent) AddRightClickEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.RightClickEvents = append(r.RightClickEvents, events...)
	return r
}

func (component MouseEventsComponent) AddDoubleRightClickEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.DoubleRightClickEvents = append(r.DoubleRightClickEvents, events...)
	return r
}

func (component MouseEventsComponent) AddMouseEnterEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.MouseEnterEvents = append(r.MouseEnterEvents, events...)
	return r
}

func (component MouseEventsComponent) AddMouseLeaveEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.MouseLeaveEvents = append(r.MouseLeaveEvents, events...)
	return r
}

func (component MouseEventsComponent) AddHoverEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.HoverEvents = append(r.HoverEvents, events...)
	return r
}

func (component MouseEventsComponent) AddDragEvents(events ...any) MouseEventsComponent {
	r := component.Clone()
	r.DragEvents = append(r.DragEvents, events...)
	return r
}

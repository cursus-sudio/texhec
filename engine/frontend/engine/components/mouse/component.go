package mouse

//

// keeps element selected even if user drags outside
type KeepSelected struct{}

//

type DragEvents struct {
	Events []any
}

func NewDragEvents(events ...any) DragEvents {
	return DragEvents{events}
}

//

type MouseEvents struct {
	LeftClickEvents        []any
	DoubleLeftClickEvents  []any
	RightClickEvents       []any
	DoubleRightClickEvents []any

	MouseEnterEvents []any
	MouseLeaveEvents []any

	// hover event is triggered every frame object is hovered
	HoverEvent []any
}

func NewMouseEvents() MouseEvents {
	return MouseEvents{}
}

func (component MouseEvents) Clone() MouseEvents {
	return MouseEvents{
		LeftClickEvents:        component.LeftClickEvents,
		DoubleLeftClickEvents:  component.DoubleLeftClickEvents,
		RightClickEvents:       component.RightClickEvents,
		DoubleRightClickEvents: component.DoubleRightClickEvents,

		MouseEnterEvents: component.MouseEnterEvents,
		MouseLeaveEvents: component.MouseLeaveEvents,

		HoverEvent: component.HoverEvent,
	}
}

func (component MouseEvents) AddLeftClickEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.LeftClickEvents = append(r.LeftClickEvents, events...)
	return r
}

func (component MouseEvents) AddDoubleLeftClickEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.DoubleLeftClickEvents = append(r.DoubleLeftClickEvents, events...)
	return r
}

func (component MouseEvents) AddRightClickEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.RightClickEvents = append(r.RightClickEvents, events...)
	return r
}

func (component MouseEvents) AddDoubleRightClickEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.DoubleRightClickEvents = append(r.DoubleRightClickEvents, events...)
	return r
}

func (component MouseEvents) AddMouseEnterEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.MouseEnterEvents = append(r.MouseEnterEvents, events...)
	return r
}

func (component MouseEvents) AddMouseLeaveEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.MouseLeaveEvents = append(r.MouseLeaveEvents, events...)
	return r
}

func (component MouseEvents) AddMouseHoverEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.HoverEvent = append(r.HoverEvent, events...)
	return r
}

package inputs

type Hovered struct{}

func NewHovered() Hovered { return Hovered{} }

//

type Dragged struct{}

func NewDragged() Dragged { return Dragged{} }

//

//

// keeps element selected even if user drags outside
type KeepSelected struct{}

//

type MouseEvents struct {
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

		HoverEvents: component.HoverEvents,
		DragEvents:  component.DragEvents,
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

func (component MouseEvents) AddHoverEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.HoverEvents = append(r.HoverEvents, events...)
	return r
}

func (component MouseEvents) AddDragEvents(events ...any) MouseEvents {
	r := component.Clone()
	r.DragEvents = append(r.DragEvents, events...)
	return r
}

package inputs

import (
	"frontend/services/media/window"
	"shared/services/ecs"
)

type System ecs.SystemRegister

type QuitEvent struct{}

func NewQuitEvent() QuitEvent { return QuitEvent{} }

// this event is called when nothing is dragged
type DragEvent struct {
	Camera   ecs.EntityID
	From, To window.MousePos // from and to is normalized
}

//

type ApplyDragEvent interface {
	Apply(DragEvent) (event any)
}

//

type SynchronizePositionEvent DragEvent

func (SynchronizePositionEvent) Apply(dragEvent DragEvent) any {
	return SynchronizePositionEvent(dragEvent)
}

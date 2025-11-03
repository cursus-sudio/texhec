package inputs

import (
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type System ecs.SystemRegister

type QuitEvent struct{}

func NewQuitEvent() QuitEvent { return QuitEvent{} }

// this event is called when nothing is dragged
type DragEvent struct {
	Camera   ecs.EntityID
	From, To mgl32.Vec2 // from and to is normalized
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

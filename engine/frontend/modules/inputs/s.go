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
	From, To mgl32.Vec2
}

// TODO
// disscuss merging inputs and render named as media

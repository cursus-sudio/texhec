package internal

import (
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"sync"

	"github.com/ogiusek/events"
)

func NewLastSystem[Component transition.Lerp[Component]]() smooth.StopSystem {
	mutex := &sync.Mutex{}
	return ecs.NewSystemRegister(func(w smooth.World) error {
		mutex.Lock()
		defer mutex.Unlock()

		t, ok := ecs.GetGlobal[tool[Component]](w)
		if !ok {
			t = NewTool[Component](w)
			w.SaveGlobal(t)
		}

		events.Listen(w.EventsBuilder(), func(tick frames.TickEvent) {
			r, ok := w.Record().Entity().Stop(t.recordingID)
			if !ok {
				return
			}
			for _, entity := range r.Entities.GetIndices() {
				beforeComponents, ok := r.Entities.Get(entity)
				if !ok || beforeComponents == nil {
					continue
				}
				before, ok := beforeComponents[0].(Component)
				if !ok {
					continue
				}
				after, ok := t.componentArray.Get(entity)
				if !ok {
					continue
				}
				lerpComponent := transition.NewTransition(before, after, tick.Delta)
				t.lerpArray.Set(entity, lerpComponent)
			}
		})

		return nil
	})
}

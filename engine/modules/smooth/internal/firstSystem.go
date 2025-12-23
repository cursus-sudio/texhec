package internal

import (
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"sync"

	"github.com/ogiusek/events"
)

func NewFirstSystem[Component transition.Lerp[Component]]() smooth.StartSystem {
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
			for _, entity := range t.lerpArray.GetEntities() {
				transitionComponent, ok := t.lerpArray.Get(entity)
				if !ok {
					continue
				}
				t.lerpArray.Remove(entity)
				t.componentArray.Set(entity, transitionComponent.To)
			}

			t.recordingID = w.Record().Entity().StartBackwardsRecording(t.config)
		})

		return nil
	})
}

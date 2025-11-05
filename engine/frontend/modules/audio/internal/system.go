package internal

import (
	"frontend/modules/audio"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type system struct {
}

func NewSystem(s Service) audio.System {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		events.ListenE(w.EventsBuilder(), func(e audio.StopEvent) error {
			return s.Stop(e.Channel)
		})
		events.ListenE(w.EventsBuilder(), func(e audio.PlayEvent) error {
			return s.Play(e.Channel, e.Asset)
		})
		events.ListenE(w.EventsBuilder(), func(e audio.QueueEvent) error {
			return s.Queue(e.Channel, e.Asset)
		})
		events.ListenE(w.EventsBuilder(), func(e audio.QueueEndlessEvent) error {
			return s.QueueEndless(e.Channel, e.Asset)
		})
		events.ListenE(w.EventsBuilder(), func(e audio.SetMasterVolumeEvent) error {
			return s.SetMasterVolume(e.Volume)
		})
		events.ListenE(w.EventsBuilder(), func(e audio.SetChannelVolumeEvent) error {
			return s.SetChannelVolume(e.Channel, e.Volume)
		})
		return nil
	})
}

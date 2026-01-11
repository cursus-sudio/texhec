package internal

import (
	"engine/modules/audio"
	"engine/services/ecs"

	"github.com/ogiusek/events"
)

func NewSystem(
	s Service,
	eventsBuilder events.Builder,
) audio.System {
	return ecs.NewSystemRegister(func() error {
		events.ListenE(eventsBuilder, func(e audio.StopEvent) error {
			return s.Stop(e.Channel)
		})
		events.ListenE(eventsBuilder, func(e audio.PlayEvent) error {
			return s.Play(e.Channel, e.Asset)
		})
		events.ListenE(eventsBuilder, func(e audio.QueueEvent) error {
			return s.Queue(e.Channel, e.Asset)
		})
		events.ListenE(eventsBuilder, func(e audio.QueueEndlessEvent) error {
			return s.QueueEndless(e.Channel, e.Asset)
		})
		events.ListenE(eventsBuilder, func(e audio.SetMasterVolumeEvent) error {
			return s.SetMasterVolume(e.Volume)
		})
		events.ListenE(eventsBuilder, func(e audio.SetChannelVolumeEvent) error {
			return s.SetChannelVolume(e.Channel, e.Volume)
		})
		return nil
	})
}

package audio

import (
	"engine/services/ecs"
)

type Channel int
type Volume float32 // volume is normalized

type Service interface {
	PlayerService
	VolumeService
}

type PlayerService interface {
	Stop(Channel) error
	Play(Channel, ecs.EntityID) error
	Queue(Channel, ecs.EntityID) error
	QueueEndless(Channel, ecs.EntityID) error
}

type VolumeService interface {
	SetMasterVolume(Volume) error
	SetChannelVolume(Channel, Volume) error
}

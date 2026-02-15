package audio

import "engine/modules/assets"

type Channel int
type Volume float32 // volume is normalized

type Service interface {
	PlayerService
	VolumeService
}

type PlayerService interface {
	Stop(Channel) error
	Play(Channel, assets.ID) error
	Queue(Channel, assets.ID) error
	QueueEndless(Channel, assets.ID) error
}

type VolumeService interface {
	SetMasterVolume(Volume) error
	SetChannelVolume(Channel, Volume) error
}

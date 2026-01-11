package audio

import "engine/services/assets"

type Channel int
type Volume float32 // volume is normalized

type Service interface {
	PlayerService
	VolumeService
}

type PlayerService interface {
	Stop(Channel) error
	Play(Channel, assets.AssetID) error
	Queue(Channel, assets.AssetID) error
	QueueEndless(Channel, assets.AssetID) error
}

type VolumeService interface {
	SetMasterVolume(Volume) error
	SetChannelVolume(Channel, Volume) error
}

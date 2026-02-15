package audio

import (
	"engine/modules/assets"
	"engine/services/ecs"
)

type System ecs.SystemRegister

type StopEvent struct {
	Channel Channel
}

func NewStopEvent(channel Channel) StopEvent {
	return StopEvent{Channel: channel}
}

//

type PlayEvent struct {
	Channel Channel
	Asset   assets.ID
}

func NewPlayEvent(channel Channel, assetID assets.ID) PlayEvent {
	return PlayEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEvent struct {
	Channel Channel
	Asset   assets.ID
}

func NewQueueEvent(channel Channel, assetID assets.ID) QueueEvent {
	return QueueEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEndlessEvent struct {
	Channel Channel
	Asset   assets.ID
}

func NewQueueEndlessEvent(channel Channel, assetID assets.ID) QueueEndlessEvent {
	return QueueEndlessEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type SetMasterVolumeEvent struct {
	Volume Volume
}

func NewSetMasterVolumeEvent(volume Volume) SetMasterVolumeEvent {
	return SetMasterVolumeEvent{Volume: volume}
}

//

type SetChannelVolumeEvent struct {
	Channel Channel
	Volume  Volume
}

func NewSetChannelVolumeEvent(channel Channel, volume Volume) SetChannelVolumeEvent {
	return SetChannelVolumeEvent{
		Channel: channel,
		Volume:  volume,
	}
}

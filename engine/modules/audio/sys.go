package audio

import (
	"engine/services/assets"
	"engine/services/ecs"
)

type System ecs.SystemRegister[ecs.World]

type StopEvent struct {
	Channel Channel
}

func NewStopEvent(channel Channel) StopEvent {
	return StopEvent{Channel: channel}
}

//

type PlayEvent struct {
	Channel Channel
	Asset   assets.AssetID
}

func NewPlayEvent(channel Channel, assetID assets.AssetID) PlayEvent {
	return PlayEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEvent struct {
	Channel Channel
	Asset   assets.AssetID
}

func NewQueueEvent(channel Channel, assetID assets.AssetID) QueueEvent {
	return QueueEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEndlessEvent struct {
	Channel Channel
	Asset   assets.AssetID
}

func NewQueueEndlessEvent(channel Channel, assetID assets.AssetID) QueueEndlessEvent {
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

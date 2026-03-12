package audio

import (
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
	Asset   ecs.EntityID
}

func NewPlayEvent(channel Channel, assetID ecs.EntityID) PlayEvent {
	return PlayEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEvent struct {
	Channel Channel
	Asset   ecs.EntityID
}

func NewQueueEvent(channel Channel, assetID ecs.EntityID) QueueEvent {
	return QueueEvent{
		Channel: channel,
		Asset:   assetID,
	}
}

//

type QueueEndlessEvent struct {
	Channel Channel
	Asset   ecs.EntityID
}

func NewQueueEndlessEvent(channel Channel, assetID ecs.EntityID) QueueEndlessEvent {
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

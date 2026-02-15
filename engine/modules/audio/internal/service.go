package internal

import (
	"engine/modules/assets"
	"engine/modules/audio"
	"engine/services/datastructures"
	"sync"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/mix"
)

type Service interface {
	audio.PlayerService
	audio.VolumeService
}

type audioService struct {
	assets assets.Service

	mutex *sync.Mutex

	masterVolume    audio.Volume
	channelsVolumes datastructures.SparseArray[audio.Channel, audio.Volume]
}

func NewService(c ioc.Dic) Service {
	return &audioService{
		assets: ioc.Get[assets.Service](c),

		mutex:           &sync.Mutex{},
		masterVolume:    1,
		channelsVolumes: datastructures.NewSparseArray[audio.Channel, audio.Volume](),
	}
}

func (s *audioService) Stop(channel audio.Channel) error {
	mix.HaltChannel(int(channel))
	return nil
}

func (s *audioService) Play(channel audio.Channel, asset assets.ID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	audioAsset, err := assets.GetAsset[audio.AudioAsset](s.assets, asset)
	if err != nil {
		return err
	}
	if _, err := audioAsset.Chunk().Play(int(channel), 0); err != nil {
		return err
	}
	return nil
}
func (s *audioService) QueueEndless(channel audio.Channel, asset assets.ID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	audioAsset, err := assets.GetAsset[audio.AudioAsset](s.assets, asset)
	if err != nil {
		return err
	}
	if _, err := audioAsset.Chunk().PlayTimed(int(channel), 0, -1); err != nil {
		return err
	}
	return nil
}
func (s *audioService) Queue(channel audio.Channel, asset assets.ID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	audioAsset, err := assets.GetAsset[audio.AudioAsset](s.assets, asset)
	if err != nil {
		return err
	}
	if _, err := audioAsset.Chunk().PlayTimed(int(channel), 1, -1); err != nil {
		return err
	}
	return nil
}

func (s *audioService) SetMasterVolume(masterVolume audio.Volume) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.masterVolume = masterVolume
	mix.Volume(-1, int(masterVolume*mix.MAX_VOLUME))
	for _, channel := range s.channelsVolumes.GetIndices() {
		volume, _ := s.channelsVolumes.Get(channel)
		mix.Volume(int(channel), int(volume*masterVolume*mix.MAX_VOLUME))
	}
	return nil
}
func (s *audioService) SetChannelVolume(channel audio.Channel, volume audio.Volume) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	mix.Volume(int(channel), int(volume*s.masterVolume*mix.MAX_VOLUME))
	return nil
}

package audio

import (
	"github.com/veandco/go-sdl2/mix"
)

type AudioAsset interface {
	Chunk() *mix.Chunk
	Release()
}

type audioAsset struct {
	chunk *mix.Chunk
}

func NewAudioAsset(chunk *mix.Chunk) AudioAsset {
	return &audioAsset{
		chunk: chunk,
	}
}

func (a *audioAsset) Chunk() *mix.Chunk { return a.chunk }
func (a *audioAsset) Release()          { a.chunk.Free() }

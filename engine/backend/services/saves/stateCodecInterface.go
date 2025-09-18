package saves

import (
	"errors"
	"sync"
)

type SaveData []byte

func NewSaveData(bytes []byte) SaveData {
	return SaveData(bytes)
}

func (save *SaveData) Bytes() []byte {
	return []byte(*save)
}

var (
	ErrInvalidSaveFormat error = errors.New("invalid save format")
)

type StateCodecRWMutex struct{ mutex *sync.RWMutex }

func newStateCodecRWMutex() StateCodecRWMutex          { return StateCodecRWMutex{mutex: &sync.RWMutex{}} }
func (mutex StateCodecRWMutex) RWMutex() *sync.RWMutex { return mutex.mutex }

type StateCodec interface {
	Serialize() (SaveData, error)
	// can return:
	// - ErrInvalidSaveFormat
	Load(SaveData) error
}

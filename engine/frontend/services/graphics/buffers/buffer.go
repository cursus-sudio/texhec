package buffers

import (
	"reflect"
	"shared/services/datastructures"
	"sort"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Buffer[Stored comparable] interface {
	ID() uint32
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
	Release()

	Flush()
}

type buffer[Stored comparable] struct {
	mutex  sync.Locker
	target uint32
	usage  uint32
	buffer uint32

	elementSize int
	datastructures.TrackingArray[Stored]
	bufferLen int
}

func NewBuffer[Stored comparable](
	target uint32, // gl.SHADER_STORAGE_BUFFER / gl.DRAW_INDIRECT_BUFFER
	usage uint32, // gl.STATIC_DRAW / gl.DYNAMIC_DRAW
	bufferID uint32,
) Buffer[Stored] {
	mutex := &sync.Mutex{}
	return &buffer[Stored]{
		mutex:  mutex,
		target: target,
		usage:  usage,
		buffer: bufferID,

		elementSize:   int(reflect.TypeFor[Stored]().Size()),
		TrackingArray: datastructures.NewThreadSafeTrackingArray[Stored](mutex),
		bufferLen:     0,
	}
}

func (s *buffer[Stored]) ID() uint32 { return s.buffer }

func (s *buffer[Stored]) CheckBufferSize() bool {
	elementsCount := len(s.TrackingArray.Get())
	if s.bufferLen != elementsCount {
		s.bufferLen = elementsCount
		return true
	}
	return false
}

func (s *buffer[Stored]) Release() {
	gl.DeleteBuffers(1, &s.buffer)
}

func (s *buffer[Stored]) Flush() {
	changes := s.TrackingArray.Changes()
	s.TrackingArray.ClearChanges()

	if len(changes) == 0 {
		return
	}

	gl.BindBuffer(s.target, s.buffer)
	defer gl.BindBuffer(s.target, 0)

	arr := s.TrackingArray.Get()
	if resized := s.CheckBufferSize(); resized {
		gl.BufferData(s.target, s.bufferLen*s.elementSize, gl.Ptr(arr), s.usage)
		return
	}

	sort.Slice(changes, func(i, j int) bool { return changes[i].Index > changes[j].Index })

	var offset int = changes[0].Index
	var size int = 1

	for _, changed := range changes[1:] {
		if changed.Index == offset+size {
			size += 1
			continue
		}
		gl.BufferSubData(s.target, offset*s.elementSize, size*s.elementSize, gl.Ptr(arr[offset:offset+size]))
		offset = changed.Index
		size = 1
	}
	gl.BufferSubData(s.target, offset*s.elementSize, size*s.elementSize, gl.Ptr(arr[offset:offset+size]))
}

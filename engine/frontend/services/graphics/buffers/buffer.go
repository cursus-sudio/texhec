package buffers

import (
	"frontend/services/datastructures"
	"reflect"
	"sort"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Buffer[Stored comparable] interface {
	ID() uint32
	Data() []Stored
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

func (s *buffer[Stored]) Data() []Stored { return s.TrackingArray.Get() }

func (s *buffer[Stored]) CheckBufferSize() bool {
	resizedBuffer := false
	if s.bufferLen == 0 {
		resizedBuffer = true
		s.bufferLen = 1
	}
	for len(s.TrackingArray.Get())*2 < s.bufferLen-1 && s.bufferLen > 1 {
		resizedBuffer = true
		s.bufferLen /= 2
	}
	for len(s.TrackingArray.Get()) > s.bufferLen {
		resizedBuffer = true
		s.bufferLen *= 2
	}
	return resizedBuffer
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

	arr := s.TrackingArray.Get()
	if resized := s.CheckBufferSize(); resized {
		gl.BindBuffer(s.target, s.buffer)
		gl.BufferData(s.target, s.bufferLen*s.elementSize, gl.Ptr(arr), s.usage)
		gl.BindBuffer(s.target, 0)
		return
	}

	gl.BindBuffer(s.target, s.buffer)

	orderedChanges := make([]int, 0, len(changes))
	for index := range changes {
		orderedChanges = append(orderedChanges, index)
	}

	sort.Ints(orderedChanges)

	var offset int = orderedChanges[0]
	var size int = 1

	for _, changed := range orderedChanges[1:] {
		if changed == offset+size {
			size += 1
			continue
		}
		gl.BufferSubData(s.target, offset*s.elementSize, size*s.elementSize, gl.Ptr(arr[offset:offset+size]))
		offset = changed
		size = 1
	}
	gl.BufferSubData(s.target, offset*s.elementSize, size*s.elementSize, gl.Ptr(arr[offset:offset+size]))
	gl.BindBuffer(s.target, 0)
}

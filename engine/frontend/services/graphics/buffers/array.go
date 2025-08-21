package buffers

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrOutOfBounds error = errors.New("out of bounds")
)

type TrackingArray[Stored comparable] interface {
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error

	Changes() []int
	ClearChanges()
}

type array[Stored comparable] struct {
	data    []Stored
	changes map[int]struct{}
}

func NewArray[Stored comparable]() TrackingArray[Stored] {
	return &array[Stored]{
		changes: map[int]struct{}{},
	}
}

func (s *array[Stored]) Get() []Stored { return s.data }

func (s *array[Stored]) Add(elements ...Stored) {
	for i := 0; i < len(elements); i++ {
		s.changes[i+len(s.data)] = struct{}{}
	}
	s.data = append(s.data, elements...)
}

func (s *array[Stored]) Set(index int, e Stored) error {
	if len(s.data) <= index {
		return ErrOutOfBounds
	}
	if s.data[index] == e {
		return nil
	}
	s.changes[index] = struct{}{}
	s.data[index] = e

	return nil
}

func (s *array[Stored]) Remove(indices ...int) error {
	for _, index := range indices {
		if index >= len(s.data) {
			return ErrOutOfBounds
		}
	}

	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		s.changes[index] = struct{}{}
		s.changes[len(s.data)-1] = struct{}{}

		s.data[index] = s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
	}
	return nil
}

func (s *array[Stored]) Changes() []int {
	changes := make([]int, 0, len(s.changes))
	for k := range s.changes {
		changes = append(changes, k)
	}
	sort.Ints(changes)
	return changes
}

func (s *array[Stored]) ClearChanges() { s.changes = make(map[int]struct{}) }

//

type threadSafeArr[Stored comparable] struct {
	mutex sync.Locker
	TrackingArray[Stored]
}

func NewThreadSafeTrackingArray[Stored comparable](mutex sync.Locker) TrackingArray[Stored] {
	return &threadSafeArr[Stored]{
		mutex:         mutex,
		TrackingArray: NewArray[Stored](),
	}
}

//	func (arr *threadSafeArr[Stored]) Get() []Stored {
//		arr.mutex.Lock()
//		defer arr.mutex.Unlock()
//		return arr.TrackingArray.Get()
//	}
func (arr *threadSafeArr[Stored]) Add(elements ...Stored) {
	arr.mutex.Lock()
	defer arr.mutex.Unlock()
	arr.TrackingArray.Add(elements...)
}
func (arr *threadSafeArr[Stored]) Set(index int, e Stored) error {
	arr.mutex.Lock()
	defer arr.mutex.Unlock()
	return arr.TrackingArray.Set(index, e)
}
func (arr *threadSafeArr[Stored]) Remove(indices ...int) error {
	arr.mutex.Lock()
	defer arr.mutex.Unlock()
	return arr.TrackingArray.Remove(indices...)
}

func (arr *threadSafeArr[Stored]) Changes() []int {
	arr.mutex.Lock()
	return arr.TrackingArray.Changes()
}

func (arr *threadSafeArr[Stored]) ClearChanges() {
	arr.TrackingArray.ClearChanges()
	arr.mutex.Unlock()
}

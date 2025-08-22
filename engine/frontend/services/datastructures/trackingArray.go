package datastructures

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

	Changes() map[int]Stored
	ClearChanges()
}

type trackingArray[Stored comparable] struct {
	data    []Stored
	changes map[int]Stored // original
}

func NewTrackingArray[Stored comparable]() TrackingArray[Stored] {
	return &trackingArray[Stored]{
		changes: map[int]Stored{},
	}
}

func (s *trackingArray[Stored]) Get() []Stored { return s.data }

func (s *trackingArray[Stored]) Add(elements ...Stored) {
	for i := 0; i < len(elements); i++ {
		var zero Stored
		s.changes[i+len(s.data)] = zero
	}
	s.data = append(s.data, elements...)
}

func (s *trackingArray[Stored]) Set(index int, e Stored) error {
	if len(s.data) <= index {
		return ErrOutOfBounds
	}
	original := s.data[index]
	if s.data[index] == e {
		return nil
	}
	if _, ok := s.changes[index]; !ok {
		s.changes[index] = original
	}
	s.data[index] = e

	return nil
}

func (s *trackingArray[Stored]) Remove(indices ...int) error {
	for _, index := range indices {
		if index >= len(s.data) {
			return ErrOutOfBounds
		}
	}

	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		indexOriginal, ok := s.changes[index]
		if !ok {
			indexOriginal = s.data[index]
		}
		lastOriginal, ok := s.changes[len(s.data)-1]
		if !ok {
			lastOriginal = s.data[len(s.data)-1]
		}
		s.changes[index], s.changes[len(s.data)-1] = indexOriginal, lastOriginal

		// s.changes[index], s.changes[len(s.data)-1] = s.data[len(s.data)-1], s.data[index]

		s.data[index] = s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
	}
	return nil
}

func (s *trackingArray[Stored]) Changes() map[int]Stored {
	return s.changes
}

func (s *trackingArray[Stored]) ClearChanges() { s.changes = map[int]Stored{} }

//

type threadSafeArr[Stored comparable] struct {
	mutex sync.Locker
	TrackingArray[Stored]
}

func NewThreadSafeTrackingArray[Stored comparable](mutex sync.Locker) TrackingArray[Stored] {
	return &threadSafeArr[Stored]{
		mutex:         mutex,
		TrackingArray: NewTrackingArray[Stored](),
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

func (arr *threadSafeArr[Stored]) Changes() map[int]Stored {
	arr.mutex.Lock()
	return arr.TrackingArray.Changes()
}

func (arr *threadSafeArr[Stored]) ClearChanges() {
	arr.TrackingArray.ClearChanges()
	arr.mutex.Unlock()
}

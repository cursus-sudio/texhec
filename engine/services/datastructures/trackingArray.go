package datastructures

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrOutOfBounds error = errors.New("out of bounds")
)

type Change[Stored comparable] struct {
	Index int
	From  *Stored
}

type TrackingArray[Stored comparable] interface {
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error

	Changes() []Change[Stored]
	ClearChanges()
}

type trackingArray[Stored comparable] struct {
	Data       []Stored
	ChangesMap map[int]*Stored
}

func NewTrackingArray[Stored comparable]() TrackingArray[Stored] {
	return &trackingArray[Stored]{
		ChangesMap: map[int]*Stored{},
	}
}

func (s *trackingArray[Stored]) Get() []Stored { return s.Data }

func (s *trackingArray[Stored]) Add(elements ...Stored) {
	for i := 0; i < len(elements); i++ {
		s.ChangesMap[i+len(s.Data)] = nil
	}
	s.Data = append(s.Data, elements...)
}

func (s *trackingArray[Stored]) Set(index int, e Stored) error {
	if len(s.Data) <= index {
		return ErrOutOfBounds
	}
	original := s.Data[index]
	if s.Data[index] == e {
		return nil
	}
	if _, ok := s.ChangesMap[index]; !ok {
		s.ChangesMap[index] = &original
	}
	s.Data[index] = e

	return nil
}

func (s *trackingArray[Stored]) Remove(indices ...int) error {
	for _, index := range indices {
		if index >= len(s.Data) {
			return ErrOutOfBounds
		}
	}

	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		if _, ok := s.ChangesMap[index]; !ok {
			e := s.Data[index]
			s.ChangesMap[index] = &e
		}
		if _, ok := s.ChangesMap[len(s.Data)-1]; !ok {
			e := s.Data[len(s.Data)-1]
			s.ChangesMap[len(s.Data)-1] = &e
		}

		s.Data[index] = s.Data[len(s.Data)-1]
		s.Data = s.Data[:len(s.Data)-1]
	}
	return nil
}

func (s *trackingArray[Stored]) Changes() []Change[Stored] {
	changes := make([]Change[Stored], 0, len(s.ChangesMap))
	for index, from := range s.ChangesMap {
		changes = append(changes, Change[Stored]{index, from})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Index > changes[j].Index })
	return changes
}

func (s *trackingArray[Stored]) ClearChanges() { s.ChangesMap = map[int]*Stored{} }

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

func (arr *threadSafeArr[Stored]) Changes() []Change[Stored] {
	arr.mutex.Lock()
	return arr.TrackingArray.Changes()
}

func (arr *threadSafeArr[Stored]) ClearChanges() {
	arr.TrackingArray.ClearChanges()
	arr.mutex.Unlock()
}

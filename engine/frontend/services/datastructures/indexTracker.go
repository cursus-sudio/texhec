package datastructures

type Set[Stored comparable] interface {
	Get() []Stored
	GetStored(index int) (element Stored, ok bool)
	GetIndex(element Stored) (index int, ok bool)
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
}

type set[Stored comparable] struct {
	TrackingArray[Stored]
	indices map[Stored]int
}

func NewSet[Stored comparable]() Set[Stored] {
	return &set[Stored]{
		TrackingArray: NewTrackingArray[Stored](),
		indices:       map[Stored]int{},
	}
}

func (s *set[Stored]) UpdateIndices() {
	changes := s.TrackingArray.Changes()
	if len(changes) == 0 {
		return
	}
	elements := s.TrackingArray.Get()
	s.TrackingArray.ClearChanges()
	for index, original := range changes {
		delete(s.indices, original)
		if index < len(elements) {
			element := elements[index]
			s.indices[element] = index
		}
	}
}

func (s *set[Stored]) GetStored(index int) (Stored, bool) {
	elements := s.TrackingArray.Get()
	if len(elements) <= index {
		var zero Stored
		return zero, false
	}
	return elements[index], true
}

func (s *set[Stored]) GetIndex(e Stored) (int, bool) {
	s.UpdateIndices()
	i, ok := s.indices[e]
	return i, ok
}

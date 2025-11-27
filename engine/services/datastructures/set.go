package datastructures

type Set[Stored comparable] interface {
	Get() []Stored
	GetStored(index int) (element Stored, ok bool)
	GetIndex(element Stored) (index int, ok bool)
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
	RemoveElements(elements ...Stored)
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
	for _, original := range changes {
		if original.From != nil {
			key := *original.From
			delete(s.indices, key)
		}
		if original.Index < len(elements) {
			element := elements[original.Index]
			s.indices[element] = original.Index
		}
	}
}

func (s *set[Stored]) GetStored(index int) (Stored, bool) {
	s.UpdateIndices()
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

func (s *set[Stored]) Add(elements ...Stored) {
	elementsToAdd := make([]Stored, 0, len(elements))
	for _, element := range elements {
		_, ok := s.GetIndex(element)
		if !ok {
			elementsToAdd = append(elementsToAdd, element)
		}
	}
	s.TrackingArray.Add(elementsToAdd...)
}

func (s *set[Stored]) RemoveElements(elements ...Stored) {
	indices := make([]int, 0, len(elements))
	for _, element := range elements {
		index, ok := s.GetIndex(element)
		if !ok {
			continue
		}
		indices = append(indices, index)
	}
	s.Remove(indices...)
}

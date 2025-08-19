package arrays

type IndexTracker[Stored comparable] interface {
	Get() []Stored
	GetStored(index int) (element Stored, ok bool)
	GetIndex(element Stored) (index int, ok bool)
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
}

type indexTracker[Stored comparable] struct {
	TrackingArray[Stored]
	elements map[int]Stored
	indices  map[Stored]int
}

func NewIndexTracker[Stored comparable](
	arr TrackingArray[Stored],
) IndexTracker[Stored] {
	return &indexTracker[Stored]{
		TrackingArray: arr,
		elements:      map[int]Stored{},
		indices:       map[Stored]int{},
	}
}

func (s *indexTracker[Stored]) Update() {
	changes := s.TrackingArray.Changes()
	s.TrackingArray.ClearChanges()
	elements := s.TrackingArray.Get()
	for _, index := range changes {
		if element, ok := s.elements[index]; ok {
			delete(s.indices, element)
			delete(s.elements, index)
		}
		if index < len(elements) {
			element := elements[index]
			s.elements[index] = element
			s.indices[element] = index
		}
	}
}

func (s *indexTracker[Stored]) GetStored(index int) (Stored, bool) {
	s.Update()
	elements := s.TrackingArray.Get()
	if len(elements) <= index {
		var zero Stored
		return zero, false
	}
	return elements[index], true
}

func (s *indexTracker[Stored]) GetIndex(e Stored) (int, bool) {
	s.Update()
	i, ok := s.indices[e]
	return i, ok
}

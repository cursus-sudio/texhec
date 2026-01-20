package datastructures_test

import (
	"engine/services/datastructures"
	"testing"
)

func TestTrackingArrRemove(t *testing.T) {
	s := datastructures.NewTrackingArray[int]()
	for i := range 10 {
		s.Add(i)
	}
	removed := []int{6}
	for _, remove := range removed {
		s.Remove(remove)
		// s.RemoveElements(remove)
	}
arr:
	for i := range 10 {
		for _, r := range removed {
			if r == i {
				continue arr
			}
		}
		for _, found := range s.Get() {
			if found == i {
				continue arr
			}
		}
		// _, ok := s.GetIndex(i)
		t.Errorf("element %d got removed\n", i)
	}
}

// func TestTrackingArrayRemove(t *testing.T) {
// 	type Structure[Stored comparable] interface {
// 		Add(...Stored)
// 		Remove(...int) error
// 		Get() []Stored
// 	}
// 	added := []int{1}
// 	removedElements := []int{1}
// 	s := datastructures.NewSet[int]()
// 	structures := []Structure[int]{
// 		datastructures.NewSet[int](),
// 		buffers.NewBuffer[int](0, 0, 0),
// 		datastructures.NewTrackingArray[int](),
// 	}
// 	for i := 0; i < 10; i++ {
// 		for _, add := range added {
// 			_, ok := s.GetIndex(add)
// 			if ok {
// 				continue
// 			}
// 			for _, structure := range structures {
// 				structure.Add(add)
// 			}
// 			s.Add(add)
// 		}
// 		for _, element := range removedElements {
// 			i, ok := s.GetIndex(element)
// 			if !ok {
// 				continue
// 			}
// 			t.Errorf("got %d for %d (value %v) when state was %v\n", i, element, s.Get(), s)
// 			for _, structure := range structures {
// 				structure.Remove(i)
// 			}
// 			s.Remove(i)
// 		}
// 		added = append(added, 5, 6, 7, 8)
// 	}
// 	for i, structure := range structures {
// 		t.Errorf("structure %d is %v\n", i, structure.Get())
// 	}
// }

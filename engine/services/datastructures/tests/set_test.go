package datastructures_test

import (
	"engine/services/datastructures"
	"testing"
)

func TestSet(t *testing.T) {
	s := datastructures.NewSet[string]()
	if len(s.Get()) != 0 {
		t.Errorf("index tracker is populated after creation")
		return
	}

	element := "expected_value"
	s.Add(element)

	if len(s.Get()) != 1 {
		t.Errorf("index tracker should have one element")
		return
	}

	index, ok := s.GetIndex(element)
	if !ok {
		t.Errorf("element isn't registered properly")
		return
	}

	if value, ok := s.GetStored(index); !ok {
		t.Errorf("element isn't registered properly")
		return
	} else if value != element {
		t.Errorf("other element is registered than expected; \"%s\" == \"%s\"\n", element, value)
		return
	}

	s.Remove(index)

	if len(s.Get()) != 0 {
		t.Errorf("expected empty array after removing last element")
		return
	}

	if _, ok := s.GetIndex(element); ok {
		t.Errorf("removed element still returns index")
		return
	}
}

func TestSetRemove(t *testing.T) {
	s := datastructures.NewSet[int]()
	for i := range 10 {
		s.Add(i)
	}
	initialLen := len(s.Get())
	removed := []int{6, 7, 8, 9}
	for _, remove := range removed {
		s.RemoveElements(remove)
	}
	if len(s.Get()) != initialLen-len(removed) {
		t.Errorf("expected elements to be of len %d but its %d\n", initialLen-len(removed), len(s.Get()))
	}
arr:
	for i := range 10 {
		for _, r := range removed {
			if r == i {
				continue arr
			}
		}
		_, ok := s.GetIndex(i)
		if !ok {
			t.Errorf("element %d got removed\n", i)
		}
	}
}

package datastructures_test

import (
	"frontend/services/datastructures"
	"testing"
)

func TestIndexTracker(t *testing.T) {
	// Get() []Stored
	// GetStored(index int) (element Stored, ok bool)
	// GetIndex(element Stored) (index int, ok bool)
	// Add(elements ...Stored)
	// Set(index int, e Stored) error
	// Remove(indices ...int) error
	m := datastructures.NewSet[string]()
	if len(m.Get()) != 0 {
		t.Errorf("index tracker is populated after creation")
		return
	}

	element := "expected_value"
	m.Add(element)

	if len(m.Get()) != 1 {
		t.Errorf("index tracker should have one element")
		return
	}

	index, ok := m.GetIndex(element)
	if !ok {
		t.Errorf("element isn't registered properly")
		return
	}

	if value, ok := m.GetStored(index); !ok {
		t.Errorf("element isn't registered properly")
		return
	} else if value != element {
		t.Errorf("other element is registered than expected; \"%s\" == \"%s\"\n", element, value)
		return
	}

	if err := m.Remove(index); err != nil {
		t.Errorf("unexpeceted error during removal \"%s\"\n", err.Error())
		return
	}

	if len(m.Get()) != 0 {
		t.Errorf("expected empty array after removing last element")
		return
	}

	if _, ok := m.GetIndex(element); ok {
		t.Errorf("removed element still returns index")
		return
	}
}

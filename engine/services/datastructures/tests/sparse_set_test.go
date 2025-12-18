package datastructures_test

import (
	"engine/services/datastructures"
	"testing"
)

func TestSpraseSet(t *testing.T) {
	set := datastructures.NewSparseSet[uint8]()
	v1, v2 := uint8(2), uint8(3)

	set.Add(v1)
	set.Add(v1)
	set.Add(v2)

	values := set.GetIndices()
	if len(values) != 2 || min(values[0], values[1]) != v1 || max(values[0], values[1]) != v2 {
		t.Errorf("sparse set has invalid values. expected [%v, %v] in any order but got %v", v1, v2, values)
	}
}
